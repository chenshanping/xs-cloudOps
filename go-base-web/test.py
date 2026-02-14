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

# 环境变量控制的上限，默认值
DEFAULT_MAX_AUTHORS = 5
DEFAULT_MAX_POEMS_PER_AUTHOR = 10
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
def fetch_poems_url(author_name, limit_per_author):
    try:
        params = {'astr': urllib.parse.quote(author_name)}  # URL编码作者名称
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
def fetch_poems_info(url):
    try:
        response = requests.get(url, headers=HEADERS)
        response.raise_for_status()  # 检查请求是否成功
        soup = BeautifulSoup(response.text, 'html.parser')

        base_info = soup.select_one('#sonsyuanwen .cont div:nth-child(2)')
        poems_name = base_info.find('h1').text  # 诗名
        author_name = base_info.find('p', attrs={'class': 'source'}).get_text(strip=True)  # 作者名
        poems_content = base_info.find('div', attrs={'class': 'contson'}).get_text(strip=True)  # 诗文内容

        # 创建诗文信息字典
        poems_info = {'诗词标题': poems_name, '诗人名称': author_name, '诗文正文': poems_content}

        # 获取更多信息
        for item in soup.select('.sons .contyishang'):
            more = item.select_one('div:nth-child(1)').attrs.get('onclick', None)
            if more:
                match = re.match(PATTERN, more)
                if match:
                    typ, id, idjm = match.groups()
                    params = {'id': id, 'idjm': idjm}
                    response = requests.get(MORE_URL.format(typ), params=params, headers=HEADERS)
                    response.raise_for_status()  # 检查请求是否成功
                    more_soup = BeautifulSoup(response.text, 'html.parser')
                    item = more_soup.select_one('.contyishang')

            title = item.find('h2').text  # 获取标题
            item.find('h2').decompose()  # 删除标题元素
            content = item.get_text(strip=True)  # 获取内容
            poems_info[title] = content  # 将标题和内容添加到诗文信息中

        return poems_info  # 返回诗文信息字典

    except Exception as e:
        print(f"从 {url} 获取诗文信息时出错: {e}")
        return {}

# 保存诗文信息到 JSON
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
    all_poems = []  # 聚合所有诗文信息
    processed_authors = 0
    total_poems = 0
    try:
        for author_name in fetch_authors_name(limit_authors=max_authors):  # 遍历作者名称
            if processed_authors >= max_authors:
                break
            processed_authors += 1
            poems_count_this_author = 0
            for poems_url in fetch_poems_url(author_name=author_name, limit_per_author=max_poems_per_author):  # 获取每位作者的诗文URL
                poems_info = fetch_poems_info(poems_url)  # 获取诗文信息
                if poems_info:  # 仅在有有效信息时保存
                    all_poems.append(poems_info)  # 收集诗文信息
                    total_poems += 1
                    poems_count_this_author += 1
                if poems_count_this_author >= max_poems_per_author:
                    break
            print(f"已处理作者: {author_name}，当步诗文数: {poems_count_this_author}")
    except KeyboardInterrupt:
        pass
    finally:
        # 输出到 JSON 文件
        save_poems_info_json(all_poems)

if __name__ == '__main__':
    main()