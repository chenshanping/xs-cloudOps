import re
import urllib.parse
import requests
import json
import os
from bs4 import BeautifulSoup

# 常量定义
BASE_URL = 'https://www.gushiwen.cn/'
AUTHORS_URL = f'{BASE_URL}authors/'  # 作者页面URL
SHIWENS_URL = f'{BASE_URL}shiwens/default.aspx'  # 诗文页面URL
MORE_URL = f'{BASE_URL}nocdn/ajax{{}}.aspx'  # 获取更多内容的URL
PATTERN = r"(\w+)Show\((\d+),'([\w\d]+)'\)"  # 正则表达式，用于匹配更多信息的脚本
HEADERS = {
    'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36 Edg/133.0.0.0'
}

# 分类来源与默认值
DEFAULT_CATEGORY = '爱国'
def get_categories():
    category_list_str = os.getenv('ENV_CATEGORY_LIST', DEFAULT_CATEGORY)
    categories = [c.strip() for c in category_list_str.split(',') if c.strip()]
    if not categories:
        categories = [DEFAULT_CATEGORY]
    return categories

# 环境变量控制的上限，默认值
DEFAULT_MAX_AUTHORS = 1
DEFAULT_MAX_POEMS_PER_AUTHOR = 5
def get_limits():
    try:
        max_authors = int(os.getenv('MAX_AUTHORS', DEFAULT_MAX_AUTHORS))
    except Exception:
        max_authors = DEFAULT_MAX_AUTHORS
    try:
        max_poems = int(os.getenv('MAX_POEMS_PER_AUTHOR', DEFAULT_MAX_POEMS_PER_AUTHOR))
    except Exception:
        max_poems = DEFAULT_MAX_POEMS_PER_AUTHOR
    return max_authors, max_poems

# 获取所有作者的名称
def fetch_authors_name(limit_authors):
    try:
        response = requests.get(AUTHORS_URL, headers=HEADERS)
        response.raise_for_status()  # 检查请求是否成功
        soup = BeautifulSoup(response.text, 'html.parser')
        count = 0
        for item in soup.select('.typecont span a'):
            if count >= limit_authors:
                break
            yield item.text  # 生成作者名称
            count += 1
    except Exception as e:
        print(f"获取作者时出错: {e}")

# 获取特定作者的诗文URL
def fetch_poems_url(author_name, limit_per_author, category_name):
    try:
        params = {'astr': urllib.parse.quote(author_name), 'tstr': category_name}  # URL编码作者名称与分类
        response = requests.get(SHIWENS_URL, params=params, headers=HEADERS)
        response.raise_for_status()  # 检查请求是否成功
        soup = BeautifulSoup(response.text, 'html.parser')
        count = 0
        for item in soup.select('.typecont span a'):
            if count >= limit_per_author:
                break
            yield BASE_URL + item.attrs['href']  # 生成诗文的完整URL
            count += 1
    except Exception as e:
        print(f"获取作者 {author_name} 的诗文URL时出错: {e}")

# 获取诗文的具体信息
def fetch_poems_info(url, category_name=None):
    try:
        response = requests.get(url, headers=HEADERS)
        response.raise_for_status()  # 检查请求是否成功
        soup = BeautifulSoup(response.text, 'html.parser')

        # 尝试定位信息区域
        base_info = soup.select_one('#sonsyuanwen .cont div:nth-child(2)')
        if not base_info:
            # 回退策略：尝试其他常见区域
            base_info = soup.select_one('.cont') or soup.find('div', class_='cont')
        poems_name = base_info.find('h1').text if base_info and base_info.find('h1') else ''
        author_name = base_info.find('p', attrs={'class': 'source'}).get_text(strip=True) if base_info and base_info.find('p', attrs={'class': 'source'}) else ''
        poems_content = base_info.find('div', attrs={'class': 'contson'}).get_text(strip=True) if base_info and base_info.find('div', attrs={'class': 'contson'}) else ''

        # 创建诗文信息字典（字段映射为英文键）
        poems_info = {
            'poems_name': poems_name.strip(),
            'author_name': author_name.strip(),
            'poems_text': poems_content.strip(),
            'translations_and_notes': '',
            'creation_background': '',
            'appreciation': '',
            'commentary_and_analysis': '',
            'brief_analysis': ''
        }

        # 试图提取译文/注释、创作背景、鉴赏、赏析、简析等字段
        # 以下为容错性提取，若页面结构不同，可能需要调整选择器
        # 译文及注释
        trans = soup.select_one('.contyishang')
        if trans:
            poems_info['translations_and_notes'] = trans.get_text(strip=True)
        # 创作背景
        bg = soup.find(text=re.compile(r'(创作背景|背景|作者意图)'))
        if bg:
            # 简单向上寻找父元素文本
            parent = bg.parent
            if parent:
                poems_info['creation_background'] = parent.get_text(strip=True)
        # 鉴赏/赏析/简析
        arts = soup.find_all(string=lambda t: t and ('鉴赏' in t or '赏析' in t or '简析' in t))
        if arts:
            # 简单拼接相关文本
            blocks = []
            for a in arts:
                blocks.append(a.strip())
            poems_info['appreciation'] = "\n".join(blocks)

        # 简单替代：如果没有获取到，上述字段保持空字符串
        return poems_info

    except Exception as e:
        print(f"从 {url} 获取诗文信息时出错: {e}")
        return {
            'poems_name': '',
            'author_name': '',
            'poems_text': '',
            'translations_and_notes': '',
            'creation_background': '',
            'appreciation': '',
            'commentary_and_analysis': '',
            'brief_analysis': ''
        }

# 输出 JSON
def save_poems_info_json(all_poems, output_path='output.json'):
    try:
        with open(output_path, 'w', encoding='utf-8') as f:
            json.dump(all_poems, f, ensure_ascii=False, indent=2)
        print(f"已将诗文信息保存到 {output_path}")
    except Exception as e:
        print(f"保存 JSON 时出错: {e}")

# 主函数
def main():
    max_authors, max_poems_per_author = get_limits()
    categories = get_categories()
    all_poems = []  # 聚合所有诗文信息

    try:
        for category in categories:  # 逐分类轮换
            print(f"开始分类: {category}")
            processed_authors = 0
            for author_name in fetch_authors_name(limit_authors=max_authors):  # 遍历作者名称
                processed_authors += 1
                if processed_authors > max_authors:
                    break
                poems_count_this_author = 0
                for poems_url in fetch_poems_url(author_name=author_name, limit_per_author=max_poems_per_author, category_name=category):  # 获取每位作者的诗文URL
                    poems_info = fetch_poems_info(poems_url, category_name=category)  # 获取诗文信息
                    if poems_info:  # 仅在有有效信息时保存
                        poems_info['category'] = category
                        all_poems.append(poems_info)  # 收集诗文信息
                        poems_count_this_author += 1
                        if poems_count_this_author >= max_poems_per_author:
                            break
            print(f"完成分类: {category}，已处理作者数: {processed_authors}，诗文总数: {len(all_poems)}")
    except KeyboardInterrupt:
        pass
    finally:
        # 输出到 JSON 文件
        save_poems_info_json(all_poems, output_path='output.json')

if __name__ == '__main__':
    main()