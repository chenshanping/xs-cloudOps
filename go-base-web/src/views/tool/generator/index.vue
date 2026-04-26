<template>
  <div class="generator-page">
    <!-- 配置列表表格 -->
    <a-card title="代码生成器">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="handleAdd">
            <PlusOutlined /> 新增配置
          </a-button>
          <a-button @click="fetchSavedConfigs" :loading="configsLoading">
            <ReloadOutlined /> 刷新
          </a-button>
        </a-space>
      </template>

      <a-table
        :columns="configTableColumns"
        :data-source="savedConfigs"
        :loading="configsLoading"
        :pagination="false"
        row-key="id"
        size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'module_name'">
            <a-space>
              <span>{{ record.module_name }}</span>
              <a-tag v-if="isModuleGenerated(record.module_name)" color="success">已生成</a-tag>
            </a-space>
          </template>
          <template v-if="column.key === 'columns_count'">
            {{ getColumnsCount(record) }}
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-button type="link" size="small" @click="handleEdit(record)">编辑</a-button>
              <a-button type="link" size="small" @click="handleCopy(record)">复制</a-button>
              <a-button type="link" size="small" @click="handleExportConfigJSON(record)">导出JSON</a-button>
              <!-- <a-button type="link" size="small" @click="handleShowERDiagram(record)"><ApartmentOutlined /> E-R图</a-button> -->
              <a-button type="link" size="small" @click="handlePreviewFromConfig(record)">预览</a-button>
              <a-button 
                type="link" 
                size="small" 
                @click="handleGenerateFromConfig(record)"
                :disabled="isModuleGenerated(record.module_name)"
              >
                生成
              </a-button>
              <a-popconfirm title="确定删除此配置?" @confirm="handleDeleteConfig(record.id)">
                <a-button type="link" size="small" danger>删除配置</a-button>
              </a-popconfirm>
              <a-popconfirm 
                v-if="isModuleGenerated(record.module_name)"
                title="确定删除已生成的代码?" 
                @confirm="handleDeleteModule(record.module_name)"
              >
                <a-button type="link" size="small" danger>删除代码</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 编辑/新增抽屉 -->
    <a-drawer
      v-model:open="drawerVisible"
      :title="drawerTitle"
      width="100%"
      :body-style="{ paddingBottom: '80px' }"
    >
      <a-tabs v-model:activeKey="activeTab">
        <!-- Tab 1: 基础配置 -->
        <a-tab-pane key="1" tab="基础配置">
          <a-form :label-col="{ span: 3 }" :wrapper-col="{ span: 20 }">
            <a-row :gutter="24">
              <a-col :span="12">
                <a-form-item label="表名" required>
                  <a-input v-model:value="config.table_name" placeholder="数据库表名，如: sys_article" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="模块名称" required>
                  <a-input v-model:value="config.module_name" placeholder="英文名称，用于路由和文件名" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="模块描述" required>
                  <a-input v-model:value="config.description" placeholder="中文描述，用于API分组" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="作者">
                  <a-input v-model:value="config.author" placeholder="作者名称" />
                </a-form-item>
              </a-col>
              <a-col :span="24">
                <a-form-item label="生成选项" :label-col="{ span: 2 }" :wrapper-col="{ span: 22 }">
                  <a-checkbox v-model:checked="config.generate_backend">生成后端代码</a-checkbox>
                  <a-checkbox v-model:checked="config.generate_frontend" style="margin-left: 16px">生成前端代码</a-checkbox>
                  <a-checkbox v-model:checked="config.generate_sql" style="margin-left: 16px">生成建表SQL</a-checkbox>
                  <a-checkbox v-model:checked="config.enable_import_export" style="margin-left: 16px">导入导出功能</a-checkbox>
                </a-form-item>
              </a-col>
              <a-col :span="24" v-if="config.generate_frontend">
                <a-form-item label="前端路径" :label-col="{ span: 2 }" :wrapper-col="{ span: 22 }">
                  <a-input v-model:value="config.frontend_path" placeholder="留空自动检测（搜索包含web的目录）" style="width: 400px" />
                  <span style="margin-left: 8px; color: #999">例如: E:\project\my-web</span>
                </a-form-item>
              </a-col>
              <a-col :span="24">
                <a-form-item label="时间字段" :label-col="{ span: 2 }" :wrapper-col="{ span: 22 }">
                  <a-checkbox v-model:checked="config.has_created_at">CreatedAt</a-checkbox>
                  <a-checkbox v-model:checked="config.has_updated_at" style="margin-left: 16px">UpdatedAt</a-checkbox>
                  <a-checkbox v-model:checked="config.has_deleted_at" style="margin-left: 16px">DeletedAt(软删除)</a-checkbox>
                  <a-checkbox v-model:checked="config.has_created_by" style="margin-left: 16px">CreatedBy(创建人)</a-checkbox>
                  <a-checkbox v-model:checked="config.has_audit" style="margin-left: 16px">Audit(审批功能)</a-checkbox>
                </a-form-item>
              </a-col>
              <a-col :span="24">
                <a-form-item label="用户关联" :label-col="{ span: 2 }" :wrapper-col="{ span: 22 }">
                  <a-checkbox v-model:checked="config.link_to_user">关联用户表</a-checkbox>
                  <a-tooltip>
                    <template #title>
                      <div>启用后生成一对一用户关联模块</div>
                      <div style="margin-top: 4px">• 自动添加 user_id 唯一索引</div>
                      <div style="margin-top: 4px">• 生成 /my 接口供当前用户使用</div>
                      <div style="margin-top: 4px">• 自动注册到用户身份注册表</div>
                      <div style="margin-top: 4px">• 适用于：医生信息、商家信息等</div>
                    </template>
                    <QuestionCircleOutlined style="margin-left: 8px; color: #999" />
                  </a-tooltip>
                </a-form-item>
              </a-col>
              <a-col :span="24" v-if="config.link_to_user">
                <a-form-item label="身份配置" :label-col="{ span: 2 }" :wrapper-col="{ span: 22 }">
                  <a-space :size="16">
                    <span>
                      <span style="color: #666; margin-right: 4px">显示名称:</span>
                      <a-input v-model:value="config.profile_name" placeholder="默认使用模块描述" style="width: 120px" />
                    </span>
                    <span style="display: flex; align-items: center">
                      <span style="color: #666; margin-right: 4px">图标:</span>
                      <a-select
                        v-model:value="config.profile_icon"
                        style="width: 180px"
                        placeholder="选择图标"
                        allow-clear
                        show-search
                        :filter-option="iconFilterOption"
                      >
                        <a-select-opt-group label="常用图标">
                          <a-select-option v-for="icon in commonIcons" :key="icon" :value="icon">
                            <SvgIcon :name="icon" style="margin-right: 8px" />{{ icon }}
                          </a-select-option>
                        </a-select-opt-group>
                        <a-select-opt-group v-if="svgIcons.length" label="自定义SVG">
                          <a-select-option v-for="icon in svgIcons" :key="icon" :value="icon">
                            <SvgIcon :name="icon" style="margin-right: 8px" />{{ icon.replace('svg:', '') }}
                          </a-select-option>
                        </a-select-opt-group>
                      </a-select>
                      <SvgIcon v-if="config.profile_icon" :name="config.profile_icon" style="margin-left: 8px; font-size: 18px" />
                    </span>
                    <span>
                      <span style="color: #666; margin-right: 4px">限定角色:</span>
                      <a-select
                        v-model:value="config.profile_role_code"
                        style="width: 150px"
                        placeholder="不限制"
                        allow-clear
                        :options="roleList.map(r => ({ label: r.name, value: r.code }))"
                      />
                    </span>
                  </a-space>
                  <a-tooltip>
                    <template #title>
                      <div>用于个人中心动态显示用户身份</div>
                      <div style="margin-top: 4px">• 显示名称: 如"医生"、"商家"</div>
                      <div style="margin-top: 4px">• 图标: 支持 Ant Design 图标和自定义 SVG</div>
                      <div style="margin-top: 4px">• 限定角色: 只有该角色的用户才能看到此身份</div>
                    </template>
                    <QuestionCircleOutlined style="margin-left: 8px; color: #999" />
                  </a-tooltip>
                </a-form-item>
              </a-col>
              <a-col :span="24" v-if="config.has_created_by">
                <a-form-item label="身份关联" :label-col="{ span: 2 }" :wrapper-col="{ span: 22 }">
                  <a-space>
                    <a-select
                      v-model:value="config.created_by_profile_table"
                      placeholder="选择创建者身份表"
                      style="width: 200px"
                      allow-clear
                      show-search
                      :filter-option="(input: string, option: any) => option.label?.toLowerCase().includes(input.toLowerCase())"
                      @change="onCreatedByProfileTableChange"
                    >
                      <a-select-option v-for="t in profileTables" :key="t.table_name" :value="t.table_name" :label="t.table_name">
                        {{ t.table_name }}
                        <span v-if="t.table_comment" style="color: #999; margin-left: 4px">({{ t.table_comment }})</span>
                      </a-select-option>
                    </a-select>
                    <a-select
                      v-if="config.created_by_profile_table"
                      v-model:value="config.created_by_profile_field"
                      placeholder="显示字段"
                      style="width: 150px"
                      :disabled="!config.created_by_profile_table"
                    >
                      <a-select-option v-for="c in getRelationColumns(config.created_by_profile_table)" :key="c.column_name" :value="c.column_name">
                        {{ c.column_name }}<span v-if="c.comment" style="color: #999">({{ c.comment }})</span>
                      </a-select-option>
                    </a-select>
                  </a-space>
                  <a-tooltip>
                    <template #title>
                      <div>关联创建者的身份信息（如医生/商家）</div>
                      <div style="margin-top: 4px">• 身份表需要含有 user_id 字段</div>
                      <div style="margin-top: 4px">• 列表中会显示创建者的身份信息</div>
                    </template>
                    <QuestionCircleOutlined style="margin-left: 8px; color: #999" />
                  </a-tooltip>
                </a-form-item>
              </a-col>
              <a-col :span="24" v-if="config.has_created_by">
                <a-form-item label="数据隔离" :label-col="{ span: 2 }" :wrapper-col="{ span: 22 }">
                  <a-checkbox v-model:checked="config.data_isolation">启用数据隔离</a-checkbox>
                  <a-tooltip v-if="config.data_isolation">
                    <template #title>
                      <div>非管理员只能看到自己创建的数据</div>
                      <div style="margin-top: 4px">管理员角色可以看到所有数据</div>
                    </template>
                    <QuestionCircleOutlined style="margin-left: 8px; color: #999" />
                  </a-tooltip>
                  <a-select
                    v-if="config.data_isolation"
                    v-model:value="selectedAdminRoles"
                    mode="multiple"
                    placeholder="选择管理员角色"
                    style="width: 300px; margin-left: 16px"
                    :options="roleList.map(r => ({ label: r.name, value: r.id }))"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="24" v-if="config.has_created_by || config.has_audit">
                <a-form-item label="前台接口" :label-col="{ span: 2 }" :wrapper-col="{ span: 22 }">
                  <a-checkbox v-model:checked="config.generate_frontend_api">生成前台用户接口</a-checkbox>
                  <a-tooltip>
                    <template #title>
                      <div>生成前台用户使用的接口（私有路由）</div>
                      <div style="margin-top: 4px">• 不做 created_by 过滤</div>
                      <div v-if="config.has_audit" style="margin-top: 4px">• 仅返回已启用且审批通过的数据</div>
                      <div style="margin-top: 4px">• 防止爬虫，需要登录</div>
                    </template>
                    <QuestionCircleOutlined style="margin-left: 8px; color: #999" />
                  </a-tooltip>
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-tab-pane>

        <!-- Tab 2: 字段配置 -->
        <a-tab-pane key="2" tab="字段配置">
          <div style="margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center;">
            <a-space>
              <a-tooltip>
                <template #title>
                  <div><b>图片/文件字段:</b></div>
                  <div>• 单个图片/文件: 存储文件ID(uint类型)</div>
                  <div>• 多个图片/文件: 存储文件ID列表(string,逗号分隔)</div>
                  <div>• 字段名自动加后缀: _file_id 或 _file_ids</div>
                  <div style="margin-top: 8px"><b>下拉框组件:</b></div>
                  <div>• 默认不配置选项: 启用(绿)/禁用(红)</div>
                  <div>• 自定义选项: 按顺序分配颜色</div>
                  <div>• 表格中使用Tag标签显示</div>
                  <div style="margin-top: 8px"><b>关联关系:</b></div>
                  <div>• 外键字段类型需要是uint</div>
                </template>
                <a-button size="small"><QuestionCircleOutlined /> 字段说明</a-button>
              </a-tooltip>
              <a-button size="small" @click="openImportSQL"><ImportOutlined /> 导入SQL</a-button>
              <a-button size="small" type="primary" @click="openImportJSON"><CodeOutlined /> 导入JSON配置</a-button>
              <a-button size="small" @click="handleExportJSON"><ExportOutlined /> 导出配置</a-button>
            </a-space>
            <a-button type="primary" @click="addColumn"><PlusOutlined /> 添加字段</a-button>
          </div>
          <a-table
            :columns="columnTableColumns"
            :data-source="config.columns"
            :pagination="false"
            size="small"
            :row-key="(_, index) => index"
            :scroll="{ x: 2000 }"
            :custom-row="customRow"
            class="draggable-table"
          >
            <template #bodyCell="{ column, record, index }">
              <template v-if="column.key === 'column_name'">
                <a-input v-model:value="record.column_name" size="small" placeholder="字段名" @change="onColumnNameChange(record)" />
              </template>
              <template v-if="column.key === 'field_type'">
                <a-select v-model:value="record.field_type" size="small" style="width: 90px" @change="onFieldTypeChange(record)">
                  <a-select-option value="string">string</a-select-option>
                  <a-select-option value="int">int</a-select-option>
                  <a-select-option value="int64">int64</a-select-option>
                  <a-select-option value="uint">uint</a-select-option>
                  <a-select-option value="float64">float64</a-select-option>
                  <a-select-option value="bool">bool</a-select-option>
                  <a-select-option value="time.Time">time.Time</a-select-option>
                </a-select>
              </template>
              <template v-if="column.key === 'db_type'">
                <a-select v-model:value="record.db_type" size="small" style="width: 90px" allow-clear placeholder="自动">
                  <a-select-option value="VARCHAR">VARCHAR</a-select-option>
                  <a-select-option value="TEXT">TEXT</a-select-option>
                  <a-select-option value="INT">INT</a-select-option>
                  <a-select-option value="BIGINT">BIGINT</a-select-option>
                  <a-select-option value="TINYINT">TINYINT</a-select-option>
                  <a-select-option value="DECIMAL">DECIMAL</a-select-option>
                  <a-select-option value="DATETIME">DATETIME</a-select-option>
                </a-select>
              </template>
              <template v-if="column.key === 'db_length'">
                <a-input-number v-model:value="record.db_length" size="small" :min="0" style="width: 60px" />
              </template>
              <template v-if="column.key === 'comment'">
                <a-input v-model:value="record.comment" size="small" style="width: 100px" />
              </template>
              <template v-if="column.key === 'related_table'">
                <template v-if="record.column_name?.endsWith('_id') && record.column_name !== 'id'">
                  <a-select 
                    v-model:value="record.related_table" 
                    size="small" 
                    style="width: 110px" 
                    placeholder="选择关联表"
                    allow-clear
                    show-search
                    :filter-option="(input: string, option: any) => option.label?.toLowerCase().includes(input.toLowerCase())"
                    @change="onColumnRelatedTableChange(record)"
                  >
                    <a-select-option v-for="t in filteredDbTables" :key="t.table_name" :value="t.table_name" :label="t.table_name">
                      {{ t.table_name }}
                    </a-select-option>
                  </a-select>
                  <a-button 
                    v-if="record.related_table" 
                    type="link" 
                    size="small" 
                    @click="openColumnRelationConfig(index)"
                    style="padding: 0 4px"
                  >
                    <SettingOutlined />
                  </a-button>
                </template>
                <span v-else style="color: #999">-</span>
              </template>
              <template v-if="column.key === 'is_required'">
                <a-checkbox v-model:checked="record.is_required" />
              </template>
              <template v-if="column.key === 'is_searchable'">
                <a-checkbox v-model:checked="record.is_searchable" />
              </template>
              <template v-if="column.key === 'search_type'">
                <a-select v-model:value="record.search_type" size="small" style="width: 70px" :disabled="!record.is_searchable">
                  <a-select-option value="eq">=</a-select-option>
                  <a-select-option value="like">LIKE</a-select-option>
                  <a-select-option value="gt">&gt;</a-select-option>
                  <a-select-option value="gte">&gt;=</a-select-option>
                  <a-select-option value="lt">&lt;</a-select-option>
                  <a-select-option value="lte">&lt;=</a-select-option>
                </a-select>
              </template>
              <template v-if="column.key === 'is_list_visible'">
                <a-checkbox v-model:checked="record.is_list_visible" />
              </template>
              <template v-if="column.key === 'is_sortable'">
                <a-checkbox v-model:checked="record.is_sortable" />
              </template>
              <template v-if="column.key === 'sort_order'">
                <a-select v-model:value="record.sort_order" size="small" style="width: 70px" :disabled="!record.is_sortable">
                  <a-select-option value="asc">ASC</a-select-option>
                  <a-select-option value="desc">DESC</a-select-option>
                </a-select>
              </template>
              <template v-if="column.key === 'is_unique'">
                <a-checkbox v-model:checked="record.is_unique" />
              </template>
              <template v-if="column.key === 'is_form_visible'">
                <a-checkbox v-model:checked="record.is_form_visible" />
              </template>
              <template v-if="column.key === 'form_type'">
                <a-select v-model:value="record.form_type" size="small" style="width: 90px" @change="onFormTypeChange(record)">
                  <a-select-option value="input">输入框</a-select-option>
                  <a-select-option value="textarea">文本域</a-select-option>
                  <a-select-option value="number">数字</a-select-option>
                  <a-select-option value="select">下拉</a-select-option>
                  <a-select-option value="switch">开关</a-select-option>
                  <a-select-option value="date">日期</a-select-option>
                  <a-select-option value="datetime">日期时间</a-select-option>
                  <a-select-option value="image">单图片</a-select-option>
                  <a-select-option value="images">多图片</a-select-option>
                  <a-select-option value="file">单文件</a-select-option>
                  <a-select-option value="files">多文件</a-select-option>
                  <a-select-option value="editor">富文本</a-select-option>
                </a-select>
              </template>
              <template v-if="column.key === 'options'">
                <a-button
                  v-if="record.form_type === 'select'"
                  type="link"
                  size="small"
                  @click="openSelectOptions(index)"
                >
                  {{ record.dict_type ? `字典:${record.dict_type}` : `选项(${record.select_options?.length || 0})` }}
                </a-button>
                <a-button
                  v-else-if="record.form_type === 'switch'"
                  type="link"
                  size="small"
                  @click="openSwitchValues(index)"
                >
                  配置
                </a-button>
                <span v-else>-</span>
              </template>
              <template v-if="column.key === 'drag'">
                <HolderOutlined 
                  class="drag-handle" 
                  style="cursor: move; color: #999" 
                  @mousedown="handleDragHandleMouseDown(index)"
                  @mouseup="handleDragHandleMouseUp"
                />
              </template>
              <template v-if="column.key === 'action'">
                <a-space :size="0">
                  <a-dropdown>
                    <a-tooltip title="插入字段">
                      <a-button type="link" size="small"><PlusCircleOutlined /></a-button>
                    </a-tooltip>
                    <template #overlay>
                      <a-menu>
                        <a-menu-item @click="insertColumnBefore(index)">在此前插入</a-menu-item>
                        <a-menu-item @click="insertColumnAfter(index)">在此后插入</a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>
                  <a-button type="link" size="small" danger @click="removeColumn(index)">删除</a-button>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-tab-pane>

        <!-- Tab 3: 一对多/多对多关联 -->
        <a-tab-pane key="3" tab="一对多/多对多">
          <a-alert type="info" show-icon style="margin-bottom: 16px">
            <template #message>
              <span><b>belongsTo</b> 关联请在「字段配置」中设置外键字段的「关联表」，此处仅配置 hasMany 和 many2many</span>
            </template>
          </a-alert>
          <div style="margin-bottom: 16px; text-align: right;">
            <a-button type="primary" @click="addRelation"><PlusOutlined /> 添加关联</a-button>
          </div>
          <a-table :columns="relationColumns" :data-source="config.relations" :pagination="false" size="small" :row-key="(_, index) => index" :scroll="{ x: 900 }">
            <template #bodyCell="{ column, record, index }">
              <template v-if="column.key === 'relation_type'">
                <a-select v-model:value="record.relation_type" size="small" style="width: 140px">
                  <a-select-option value="hasMany">一对多(他表有外键)</a-select-option>
                  <a-select-option value="many2many">多对多</a-select-option>
                </a-select>
              </template>
              <template v-if="column.key === 'related_table'">
                <a-select 
                  v-model:value="record.related_table" 
                  size="small" 
                  style="width: 160px" 
                  placeholder="选择关联表"
                  show-search
                  :filter-option="(input: string, option: any) => option.label?.toLowerCase().includes(input.toLowerCase())"
                  @change="onRelatedTableChange(record)"
                >
                  <a-select-option v-for="t in filteredDbTables" :key="t.table_name" :value="t.table_name" :label="t.table_name">
                    {{ t.table_name }}
                    <span v-if="t.table_comment" style="color: #999; margin-left: 4px">({{ t.table_comment }})</span>
                  </a-select-option>
                </a-select>
              </template>
              <template v-if="column.key === 'related_module'">
                <a-tooltip>
                  <template #title>关联模块的 API 文件名，留空则使用表名</template>
                  <a-input v-model:value="record.related_module" size="small" placeholder="留空=表名" style="width: 110px" />
                </a-tooltip>
              </template>
              <template v-if="column.key === 'foreign_key'">
                <a-input v-model:value="record.foreign_key" size="small" placeholder="外键字段" />
              </template>
              <template v-if="column.key === 'display_field'">
                <a-select 
                  v-model:value="record.display_field" 
                  size="small" 
                  style="width: 120px" 
                  placeholder="选择"
                  :disabled="!record.related_table"
                  @change="onDisplayFieldChange(record)"
                >
                  <a-select-option v-for="c in getRelationColumns(record.related_table)" :key="c.column_name" :value="c.column_name">
                    {{ c.column_name }}<span v-if="c.comment" style="color: #999">({{ c.comment }})</span>
                  </a-select-option>
                </a-select>
              </template>
              <template v-if="column.key === 'comment'">
                <a-input v-model:value="record.comment" size="small" placeholder="自动获取" />
              </template>
              <template v-if="column.key === 'is_required'">
                <a-checkbox v-model:checked="record.is_required" :disabled="record.relation_type === 'hasMany'" />
              </template>
              <template v-if="column.key === 'join_table'">
                <a-input v-model:value="record.join_table" size="small" placeholder="中间表" :disabled="record.relation_type !== 'many2many'" />
              </template>
              <template v-if="column.key === 'use_options_api'">
                <a-space>
                  <a-tooltip>
                    <template #title>
                      <div><b>开启</b>：调用关联表的 options 接口</div>
                      <div style="margin-top: 4px">返回 {id, name, count}</div>
                      <div style="margin-top: 6px"><b>关闭</b>：调用关联表的分页列表接口</div>
                      <div style="margin-top: 6px; color: #faad14">⚠️ 关联表需要有 options 接口</div>
                    </template>
                    <a-checkbox v-model:checked="record.use_options_api" :disabled="record.relation_type === 'hasMany'" />
                  </a-tooltip>
                  <a-button 
                    v-if="record.use_options_api && record.related_table" 
                    type="link" 
                    size="small" 
                    @click="showOptionsCodePreview(record)"
                  >
                    查看代码
                  </a-button>
                </a-space>
              </template>
              <template v-if="column.key === 'action'">
                <a-button type="link" size="small" danger @click="removeRelation(index)">删除</a-button>
              </template>
            </template>
          </a-table>
        </a-tab-pane>

        <!-- Tab 4: 菜单配置 -->
        <a-tab-pane v-if="config.generate_frontend" key="4" tab="菜单配置">
          <a-form :label-col="{ span: 4 }" :wrapper-col="{ span: 16 }">
            <a-form-item label="父菜单">
              <a-tree-select
                v-model:value="menuConfig.parent_id"
                :tree-data="menuTree"
                placeholder="请选择父菜单"
                allow-clear
                tree-default-expand-all
                :field-names="{ label: 'name', value: 'id', children: 'children' }"
              />
            </a-form-item>
            <a-form-item label="菜单名称">
              <a-input v-model:value="menuConfig.menu_name" placeholder="菜单显示名称" />
            </a-form-item>
            <a-form-item label="菜单图标">
              <IconSelect v-model="menuConfig.menu_icon" />
            </a-form-item>
            <a-form-item label="菜单排序">
              <a-input-number v-model:value="menuConfig.menu_sort" :min="0" />
            </a-form-item>
            <a-form-item label="权限标识">
              <a-input :value="autoPermission" disabled placeholder="自动生成" />
              <div style="margin-top: 4px; color: #999; font-size: 12px">
                根据父菜单和模块名自动生成：{{ autoPermission || '请先选择父菜单和填写模块名称' }}
              </div>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- Tab 5: 统计配置 -->
        <a-tab-pane key="5" tab="📊 统计图表">
          <a-alert type="info" show-icon style="margin-bottom: 16px">
            <template #message>配置后将自动生成统计接口和 ECharts 图表组件，每个分组字段对应一个图表</template>
          </a-alert>
          <a-form :label-col="{ span: 4 }" :wrapper-col="{ span: 18 }">
            <a-form-item label="启用统计">
              <a-switch v-model:checked="statsConfig.enabled" />
            </a-form-item>
            <template v-if="statsConfig.enabled">
              <a-form-item label="分组图表">
                <div style="margin-bottom: 8px">
                  <a-button type="dashed" size="small" @click="addStatsChart"><PlusOutlined /> 添加图表</a-button>
                </div>
                <a-table 
                  :columns="statsChartColumns" 
                  :data-source="statsConfig.charts" 
                  :pagination="false" 
                  size="small" 
                  :row-key="(_, index) => index"
                >
                  <template #bodyCell="{ column, record, index }">
                    <template v-if="column.key === 'field'">
                      <a-select 
                        v-model:value="record.field" 
                        size="small" 
                        style="width: 160px" 
                        placeholder="选择字段"
                        @change="onStatsChartFieldChange(record)"
                      >
                        <a-select-option 
                          v-for="col in statsGroupableColumns" 
                          :key="col.column_name" 
                          :value="col.column_name"
                        >
                          {{ col.column_name }}
                          <span v-if="col.comment" style="color: #999"> ({{ col.comment }})</span>
                        </a-select-option>
                      </a-select>
                    </template>
                    <template v-if="column.key === 'chart_type'">
                      <a-select v-model:value="record.chart_type" size="small" style="width: 100px">
                        <a-select-option value="pie">🥧 饼图</a-select-option>
                        <a-select-option value="bar">📊 柱状图</a-select-option>
                      </a-select>
                    </template>
                    <template v-if="column.key === 'title'">
                      <a-input v-model:value="record.title" size="small" placeholder="自动生成" style="width: 120px" />
                    </template>
                    <template v-if="column.key === 'action'">
                      <a-button type="link" size="small" danger @click="statsConfig.charts.splice(index, 1)">删除</a-button>
                    </template>
                  </template>
                </a-table>
              </a-form-item>
              <a-form-item label="趋势图">
                <a-select
                  v-model:value="statsConfig.time_field"
                  placeholder="可选，选择时间字段生成折线趋势图"
                  allow-clear
                  style="width: 300px"
                >
                  <a-select-option value="created_at">created_at (创建时间)</a-select-option>
                  <a-select-option value="updated_at">updated_at (更新时间)</a-select-option>
                  <a-select-option 
                    v-for="col in statsTimeColumns" 
                    :key="col.column_name" 
                    :value="col.column_name"
                  >
                    {{ col.column_name }}
                    <span v-if="col.comment" style="color: #999"> ({{ col.comment }})</span>
                  </a-select-option>
                </a-select>
                <span style="margin-left: 8px; color: #999; font-size: 12px">
                  选择后生成时间趋势折线图
                </span>
              </a-form-item>
              <a-form-item label="生成预览">
                <a-card size="small" style="background: #fafafa">
                  <div style="color: #666; font-size: 13px">
                    <div v-for="(chart, i) in statsConfig.charts" :key="i" style="margin-bottom: 4px">
                      <strong>分组统计{{ i + 1 }}:</strong> 
                      <code>GET /{{ config.module_name }}/stats/{{ chart.field }}</code>
                      <span style="margin-left: 8px">按 {{ chart.field }} {{ chart.chart_type === 'pie' ? '饼图' : '柱状图' }}</span>
                    </div>
                    <div v-if="statsConfig.time_field" style="margin-top: 8px">
                      <strong>时间趋势:</strong> 
                      <code>GET /{{ config.module_name }}/stats/trend</code>
                      <span style="margin-left: 8px">按 {{ statsConfig.time_field }} 折线图</span>
                    </div>
                    <div v-if="statsConfig.charts.length === 0 && !statsConfig.time_field" style="color: #999">
                      请添加至少一个分组图表或选择时间字段
                    </div>
                  </div>
                </a-card>
              </a-form-item>
            </template>
          </a-form>
        </a-tab-pane>
      </a-tabs>

      <!-- 抽屉底部按钮 -->
      <div class="drawer-footer">
        <a-space>
          <a-button @click="drawerVisible = false">取消</a-button>
          <a-button @click="handleSaveConfig" :loading="saveLoading">保存配置</a-button>
          <a-button @click="handlePreview" :loading="previewLoading" :disabled="!config.table_name">预览代码</a-button>
          <a-button 
            v-if="isModuleGenerated(config.module_name)" 
            @click="handleShowChangeGuide" 
            :disabled="!config.table_name"
          >
            <DiffOutlined /> 变更指南
          </a-button>
          <a-button type="primary" @click="handleGenerate" :loading="generateLoading" :disabled="!config.table_name">
            保存并生成代码
          </a-button>
        </a-space>
      </div>
    </a-drawer>

    <!-- 预览弹窗 -->
    <a-modal v-model:open="previewVisible" title="代码预览" width="85%" :footer="null">
      <a-tabs v-model:activeKey="previewTab" tab-position="left">
        <a-tab-pane v-for="file in previewFiles" :key="file.path" :tab="getDisplayPath(file.path)">
          <div class="code-preview">
            <div class="code-toolbar">
              <a-space>
                <a-button v-if="file.type === 'sql'" type="primary" size="small" @click="handleExecuteSQL(file.content)" :loading="executeSqlLoading">
                  执行建SQL
                </a-button>
                <a-button size="small" @click="handleCopyCode(file.content)">
                  <CopyOutlined /> 复制代码
                </a-button>
              </a-space>
            </div>
            <pre class="code-block" :class="getCodeLanguage(file.path)"><code>{{ file.content }}</code></pre>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-modal>

    <!-- 下拉选项配置弹窗 -->
    <a-modal v-model:open="selectOptionsVisible" title="下拉选项配置" @ok="saveSelectOptions" width="600px">
      <!-- 选项来源切换 -->
      <a-radio-group v-model:value="useDictMode" style="margin-bottom: 16px">
        <a-radio-button :value="false">手动配置选项</a-radio-button>
        <a-radio-button :value="true">使用数据字典</a-radio-button>
      </a-radio-group>

      <!-- 数据字典模式 -->
      <template v-if="useDictMode">
        <a-form-item label="字典类型" :label-col="{ span: 4 }" :wrapper-col="{ span: 20 }">
          <a-select
            v-model:value="editingDictType"
            placeholder="请选择字典类型"
            show-search
            :filter-option="(input: string, option: any) => option.label?.toLowerCase().includes(input.toLowerCase())"
            style="width: 100%"
          >
            <a-select-option v-for="dt in dictTypeList" :key="dt.type" :value="dt.type" :label="dt.name">
              {{ dt.name }} <span style="color: #999">({{ dt.type }})</span>
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-alert type="info" show-icon>
          <template #message>使用数据字典</template>
          <template #description>
            <div>选项将从数据字典动态获取，便于统一管理</div>
            <div style="margin-top: 4px">可在「系统管理 → 字典管理」中维护字典数据</div>
          </template>
        </a-alert>
      </template>

      <!-- 手动配置模式 -->
      <template v-else>
        <a-alert type="info" show-icon style="margin-bottom: 16px">
          <template #message>提示</template>
          <template #description>
            <div>不配置选项则默认为启用/禁用（绿色/红色）</div>
            <div style="margin-top: 4px">自定义选项颜色顺序: 
              <a-tag color="blue">蓝</a-tag>
              <a-tag color="green">绿</a-tag>
              <a-tag color="orange">橙</a-tag>
              <a-tag color="purple">紫</a-tag>
              <a-tag color="cyan">青</a-tag>
              <a-tag color="magenta">洋红</a-tag>
              <a-tag color="gold">金</a-tag>
              <a-tag color="lime">青柠</a-tag>
            </div>
          </template>
        </a-alert>
        <a-table :columns="selectOptionColumns" :data-source="editingSelectOptions" :pagination="false" size="small">
          <template #bodyCell="{ column, record, index }">
            <template v-if="column.key === 'label'">
              <a-input v-model:value="record.label" size="small" placeholder="显示文本" />
            </template>
            <template v-if="column.key === 'value'">
              <a-input v-model:value="record.value" size="small" placeholder="值" />
            </template>
            <template v-if="column.key === 'action'">
              <a-button type="link" size="small" danger @click="editingSelectOptions.splice(index, 1)">删除</a-button>
            </template>
          </template>
        </a-table>
        <a-button type="dashed" block style="margin-top: 8px" @click="editingSelectOptions.push({ label: '', value: '' })">
          添加选项
        </a-button>
      </template>
    </a-modal>

    <!-- 开关值配置弹窗 -->
    <a-modal v-model:open="switchValuesVisible" title="开关值配置" @ok="saveSwitchValues">
      <a-form :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }">
        <a-form-item label="开启值">
          <a-input v-model:value="editingSwitchValues.active_value" placeholder="如: 1 或 true" />
        </a-form-item>
        <a-form-item label="开启文本">
          <a-input v-model:value="editingSwitchValues.active_text" placeholder="如: 启用" />
        </a-form-item>
        <a-form-item label="关闭值">
          <a-input v-model:value="editingSwitchValues.inactive_value" placeholder="如: 0 或 false" />
        </a-form-item>
        <a-form-item label="关闭文本">
          <a-input v-model:value="editingSwitchValues.inactive_text" placeholder="如: 禁用" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Options代码预览弹窗 -->
    <a-modal v-model:open="optionsCodeVisible" :title="`关联表 ${optionsCodeTable} 需要添加的代码`" width="800px" :footer="null">
      <a-alert type="warning" show-icon style="margin-bottom: 16px">
        <template #message>如果关联表已生成，需要手动添加以下代码，或重新生成关联表模块</template>
      </a-alert>
      <a-tabs>
        <a-tab-pane key="service" tab="Service">
          <div class="code-toolbar" style="margin-bottom: 8px">
            <a-button size="small" @click="handleCopyCode(optionsCodeService)"><CopyOutlined /> 复制</a-button>
          </div>
          <pre class="code-block"><code>{{ optionsCodeService }}</code></pre>
        </a-tab-pane>
        <a-tab-pane key="api" tab="API">
          <div class="code-toolbar" style="margin-bottom: 8px">
            <a-button size="small" @click="handleCopyCode(optionsCodeApi)"><CopyOutlined /> 复制</a-button>
          </div>
          <pre class="code-block"><code>{{ optionsCodeApi }}</code></pre>
        </a-tab-pane>
        <a-tab-pane key="router" tab="Router">
          <div class="code-toolbar" style="margin-bottom: 8px">
            <a-button size="small" @click="handleCopyCode(optionsCodeRouter)"><CopyOutlined /> 复制</a-button>
          </div>
          <pre class="code-block"><code>{{ optionsCodeRouter }}</code></pre>
        </a-tab-pane>
        <a-tab-pane key="frontend" tab="前端 API">
          <div class="code-toolbar" style="margin-bottom: 8px">
            <a-button size="small" @click="handleCopyCode(optionsCodeFrontend)"><CopyOutlined /> 复制</a-button>
          </div>
          <pre class="code-block"><code>{{ optionsCodeFrontend }}</code></pre>
        </a-tab-pane>
      </a-tabs>
    </a-modal>

    <!-- SQL导入弹窗 -->
    <a-modal v-model:open="importSQLVisible" title="导入SQL" width="700px" @ok="handleImportSQL">
      <a-alert type="info" show-icon style="margin-bottom: 16px">
        <template #message>支持导入CREATE TABLE语句，自动解析字段信息</template>
        <template #description>
          <div>• 自动提取表名、字段名、类型、注释</div>
          <div>• 自动转换为Go类型和表单组件</div>
          <div>• 将替换当前字段配置</div>
        </template>
      </a-alert>
      <a-textarea
        v-model:value="importSQLText"
        placeholder="粘贴CREATE TABLE语句..."
        :rows="15"
        style="font-family: monospace"
      />
    </a-modal>

    <!-- JSON配置导入弹窗 -->
    <a-modal v-model:open="importJSONVisible" title="导入JSON配置" width="800px" @ok="handleImportJSON">
      <a-alert type="info" show-icon style="margin-bottom: 16px">
        <template #message>快速导入预设的JSON配置</template>
        <template #description>
          <div>• 支持完整的GeneratorConfig JSON格式</div>
          <div>• 将替换当前所有配置（基础配置、字段、关联、菜单）</div>
          <div>• 可用于快速复制和迁移配置</div>
        </template>
      </a-alert>
      <a-textarea
        v-model:value="importJSONText"
        placeholder='粘贴JSON配置，例如：
{
  "table_name": "livestock_category",
  "module_name": "livestock_category",
  "description": "牲畜分类",
  "generate_backend": true,
  "generate_frontend": true,
  ...
}'
        :rows="18"
        style="font-family: monospace; font-size: 12px"
      />
    </a-modal>

    <!-- E-R图抽屉 -->
    <a-drawer
      v-model:open="erDiagramVisible"
      title="📊 概念模型 E-R 图"
      placement="bottom"
      :height="'85vh'"
    >
      <ERDiagram :entities="erEntities" :relations="erRelations" style="height: calc(85vh - 100px)" />
    </a-drawer>

    <!-- 变更指南弹窗 -->
    <a-modal v-model:open="changeGuideVisible" title="变更指南" width="1000px" :footer="null">
      <a-tabs v-model:activeKey="changeGuideTab" tab-position="left" style="min-height: 400px">
        <a-tab-pane v-for="cat in changeCategories" :key="cat.key">
          <template #tab>
            <span>{{ cat.label }}</span>
          </template>
          <!-- 变更摘要 -->
          <div v-if="cat.summary.length > 0" style="margin-bottom: 16px; padding: 12px; background: #f6ffed; border: 1px solid #b7eb8f; border-radius: 4px">
            <div v-for="(item, i) in cat.summary" :key="i" style="margin-bottom: 4px">
              <CheckCircleOutlined style="color: #52c41a; margin-right: 8px" />{{ item }}
            </div>
          </div>
          <!-- 代码指南列表 -->
          <div v-for="(guide, index) in cat.guides" :key="index" style="margin-bottom: 20px; padding-bottom: 20px; border-bottom: 1px dashed #e8e8e8">
            <div style="font-weight: 500; font-size: 14px; margin-bottom: 8px; color: #1890ff">{{ guide.title }}</div>
            <div style="margin-bottom: 8px; color: #666">{{ guide.description }}</div>
            <div style="margin-bottom: 8px">
              <strong>文件：</strong>
              <code style="background: #f5f5f5; padding: 2px 6px; border-radius: 3px">{{ guide.file }}</code>
            </div>
            <div class="code-toolbar" style="margin-bottom: 8px">
              <a-button size="small" @click="handleCopyCode(guide.code)"><CopyOutlined /> 复制代码</a-button>
            </div>
            <pre class="code-block" style="max-height: 250px"><code>{{ guide.code }}</code></pre>
          </div>
          <a-empty v-if="cat.guides.length === 0 && cat.key !== 'summary'" description="无变更" />
        </a-tab-pane>
      </a-tabs>
      <a-empty v-if="changeCategories.length === 0" description="配置无变化或无法检测差异" />
    </a-modal>

    <!-- 字段关联配置弹窗 -->
    <a-modal v-model:open="columnRelationConfigVisible" title="关联配置" @ok="saveColumnRelationConfig" width="500px">
      <a-form :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }">
        <a-form-item label="关联表">
          <span>{{ editingColumnRelationIndex >= 0 ? config.columns[editingColumnRelationIndex]?.related_table : '' }}</span>
        </a-form-item>
        <a-form-item label="模块名">
          <a-input v-model:value="editingColumnRelation.related_module" placeholder="留空则使用表名" />
          <div style="color: #999; font-size: 12px; margin-top: 4px">用于 API 导入路径，如 @/api/xxx.ts</div>
        </a-form-item>
        <a-form-item label="显示字段">
          <a-select v-model:value="editingColumnRelation.display_field" placeholder="选择显示字段">
            <a-select-option v-for="c in editingColumnRelationColumns" :key="c.column_name" :value="c.column_name">
              {{ c.column_name }}<span v-if="c.comment" style="color: #999"> ({{ c.comment }})</span>
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="轻量接口">
          <a-checkbox v-model:checked="editingColumnRelation.use_options_api">
            使用 options 接口（返回 id, name, count）
          </a-checkbox>
        </a-form-item>
        <a-form-item label="左树右表">
          <a-checkbox v-model:checked="editingColumnRelation.use_tree_layout">
            启用左树右表布局
          </a-checkbox>
          <div style="color: #999; font-size: 12px; margin-top: 4px">左侧显示分类树，右侧显示表格</div>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, ReloadOutlined, QuestionCircleOutlined, CopyOutlined, ImportOutlined, CodeOutlined, ExportOutlined, PlusCircleOutlined, HolderOutlined, DiffOutlined, CheckCircleOutlined, ApartmentOutlined, SettingOutlined } from '@ant-design/icons-vue'
import {
  previewCode,
  generateCode,
  getGeneratedModules,
  deleteModule,
  saveConfig as saveConfigApi,
  getSavedConfigs,
  deleteSavedConfig,
  executeSQL,
  getTables,
  getTableColumns,
  type ColumnConfig,
  type GeneratorConfig,
  type GeneratedFile,
  type SavedConfig,
  type SelectOption,
  type SwitchValue,
  type TableInfo
} from '@/api/generator'
import { getMenuTree } from '@/api/menu'
import { getRoleList } from '@/api/role'
import { getAllDictTypes, type DictType } from '@/api/dict'
import IconSelect from '@/components/IconSelect.vue'
import SvgIcon from '@/components/SvgIcon.vue'
import ERDiagram from '@/components/ERDiagram.vue'
import type { Menu, Role } from '@/types'

// 常用 Ant Design 图标列表
const commonIcons = [
  'UserOutlined', 'TeamOutlined', 'SolutionOutlined', 'IdcardOutlined',
  'MedicineBoxOutlined', 'ShopOutlined', 'BankOutlined', 'HomeOutlined',
  'CarOutlined', 'CoffeeOutlined', 'TrophyOutlined', 'CrownOutlined',
  'StarOutlined', 'HeartOutlined', 'SmileOutlined', 'SafetyCertificateOutlined',
  'ToolOutlined', 'SettingOutlined', 'AppstoreOutlined', 'ProfileOutlined'
]

// 加载自定义 SVG 图标
const svgModules = import.meta.glob<string>('@/assets/icons/*.svg', { eager: true, query: '?raw', import: 'default' })
const svgIcons = Object.keys(svgModules).map(path => {
  const name = path.split('/').pop()?.replace('.svg', '') || ''
  return `svg:${name}`
})

// 图标搜索过滤
const iconFilterOption = (input: string, option: any) => {
  const value = option.value?.toLowerCase() || ''
  return value.includes(input.toLowerCase())
}

// 状态
const activeTab = ref('1')
const drawerVisible = ref(false)
const drawerTitle = ref('新增配置')
const menuTree = ref<Menu[]>([])
const previewVisible = ref(false)
const previewTab = ref('')
const previewFiles = ref<GeneratedFile[]>([])
const previewLoading = ref(false)
const generateLoading = ref(false)
const saveLoading = ref(false)
const executeSqlLoading = ref(false)
const generatedModules = ref<string[]>([])
const savedConfigs = ref<SavedConfig[]>([])
const configsLoading = ref(false)

// 下拉选项编辑
const selectOptionsVisible = ref(false)
const editingColumnIndex = ref(-1)
const editingSelectOptions = ref<SelectOption[]>([])
const editingDictType = ref('')  // 编辑中的字典类型
const useDictMode = ref(false)   // 是否使用字典模式

// 开关值编辑
const switchValuesVisible = ref(false)
const editingSwitchValues = reactive<SwitchValue>({
  active_value: '1',
  inactive_value: '0',
  active_text: '启用',
  inactive_text: '禁用'
})

// SQL导入
const importSQLVisible = ref(false)
const importSQLText = ref('')

// JSON配置导入
const importJSONVisible = ref(false)
const importJSONText = ref('')

// Options代码预览
const optionsCodeVisible = ref(false)
const optionsCodeTable = ref('')
const optionsCodeService = ref('')
const optionsCodeApi = ref('')
const optionsCodeRouter = ref('')
const optionsCodeFrontend = ref('')

// E-R图
const erDiagramVisible = ref(false)
interface EREntity {
  name: string
  comment: string
  columns: { name: string; comment: string; isPrimary?: boolean }[]
}
interface ERRelation {
  name: string
  from: string
  to: string
  fromCardinality: string
  toCardinality: string
}
const erEntities = ref<EREntity[]>([])
const erRelations = ref<ERRelation[]>([])

// 变更指南
const changeGuideVisible = ref(false)
const changeGuideTab = ref('summary')
interface ChangeGuide {
  title: string
  description: string
  file: string
  code: string
}
interface ChangeCategory {
  key: string
  label: string
  summary: string[]
  guides: ChangeGuide[]
}
const changeCategories = ref<ChangeCategory[]>([])

// 数据库表列表（用于关联选择）
const dbTables = ref<TableInfo[]>([])
// 过滤掉sys_开头的系统表
const filteredDbTables = computed(() => dbTables.value.filter(t => !t.table_name.startsWith('sys_')))
// 身份表列表（从已保存的配置中获取 link_to_user 为 true 的表，或含有 user_id 字段的表）
const profileTables = computed(() => {
  // 从已保存配置中获取身份表
  const profileTableNames = new Set<string>()
  for (const cfg of savedConfigs.value) {
    try {
      const parsed = JSON.parse(cfg.config_json)
      if (parsed.link_to_user) {
        profileTableNames.add(cfg.table_name)
      }
    } catch {
      // ignore
    }
  }
  // 过滤表列表，返回身份表
  return filteredDbTables.value.filter(t => profileTableNames.has(t.table_name))
})
// 关联表字段列表缓存 { tableName: columns[] }
const relationTableColumns = ref<Record<string, ColumnConfig[]>>({})
// 角色列表（用于数据隔离配置）
const roleList = ref<Role[]>([])
// 字典类型列表（用于下拉框配置）
const dictTypeList = ref<DictType[]>([])

// 配置列表表格列
const configTableColumns = [
  { title: '模块名称', key: 'module_name', dataIndex: 'module_name', width: 200 },
  { title: '描述', dataIndex: 'description', width: 200 },
  { title: '表名', dataIndex: 'table_name', width: 180 },
  { title: '字段数', key: 'columns_count', width: 80 },
  { title: '操作', key: 'action', width: 320 }
]

// 创建空字段
const createEmptyColumn = (): ColumnConfig => ({
  column_name: '',
  field_name: '',
  field_type: 'string',
  json_name: '',
  ts_type: 'string',
  comment: '',
  db_type: '',
  db_length: 0,
  default_value: '',
  is_primary_key: false,
  is_required: false,
  is_searchable: false,
  search_type: 'eq',
  is_list_visible: true,
  is_form_visible: true,
  is_sortable: false,
  sort_order: 'asc',
  is_unique: false,
  form_type: 'input',
  dict_type: '',
  select_options: [],
  switch_values: null,
  // belongsTo 关联配置
  related_table: '',
  related_module: '',
  display_field: 'name',
  use_options_api: true,
  use_tree_layout: false
})

// 配置数据
const config = reactive<GeneratorConfig>({
  id: undefined,
  table_name: '',
  module_name: '',
  description: '',
  author: '',
  generate_backend: true,
  generate_frontend: true,
  generate_sql: true,
  frontend_path: '',
  has_created_at: true,
  has_updated_at: true,
  has_deleted_at: false,
  has_created_by: false,
  created_by_profile_table: '',
  created_by_profile_field: '',
  data_isolation: false,
  admin_role_ids: '',
  has_audit: false,
  generate_frontend_api: false,
  link_to_user: false,
  profile_name: '',
  profile_icon: '',
  profile_role_code: '',
  enable_import_export: false,
  columns: [],
  relations: [],
  menu_config: null
})

const menuConfig = reactive({
  parent_id: 0,
  menu_name: '',
  menu_icon: '',
  menu_sort: 0,
  permission: ''
})

// 统计配置
const statsConfig = reactive<{
  enabled: boolean
  charts: { field: string; chart_type: string; title: string }[]
  time_field: string
}>({
  enabled: false,
  charts: [],
  time_field: ''
})

// 可用于分组的字段（外键、状态、类型等）
const statsGroupableColumns = computed(() => {
  return config.columns.filter(col => 
    col.column_name.endsWith('_id') || 
    col.column_name === 'status' || 
    col.column_name === 'type' ||
    col.form_type === 'select'
  )
})

// 时间类型字段
const statsTimeColumns = computed(() => {
  return config.columns.filter(col => 
    col.field_type === 'time.Time' || 
    col.form_type === 'date' || 
    col.form_type === 'datetime'
  )
})

// 统计图表表格列
const statsChartColumns = [
  { title: '分组字段', key: 'field', width: 180 },
  { title: '图表类型', key: 'chart_type', width: 120 },
  { title: '标题', key: 'title', width: 140 },
  { title: '操作', key: 'action', width: 80 }
]

// 添加统计图表
const addStatsChart = () => {
  statsConfig.charts.push({ field: '', chart_type: 'pie', title: '' })
}

// 图表字段变更时自动生成标题
const onStatsChartFieldChange = (chart: { field: string; chart_type: string; title: string }) => {
  if (chart.field && !chart.title) {
    const col = config.columns.find(c => c.column_name === chart.field)
    chart.title = col?.comment || chart.field
  }
}

// 当启用 link_to_user 时，自动将菜单父级设为 system(1)
watch(() => config.link_to_user, (newVal) => {
  if (newVal && menuConfig.parent_id === 0) {
    menuConfig.parent_id = 1
  }
})

// 创建者身份表选择变化时加载字段列表
const onCreatedByProfileTableChange = async (tableName: string) => {
  if (!tableName) {
    config.created_by_profile_field = ''
    return
  }
  // 加载表字段
  if (!relationTableColumns.value[tableName]) {
    try {
      const res = await getTableColumns(tableName)
      relationTableColumns.value[tableName] = res.data || []
    } catch {
      // ignore
    }
  }
  // 默认选择 name 字段
  const cols = relationTableColumns.value[tableName] || []
  const nameCol = cols.find(c => c.column_name === 'name')
  if (nameCol) {
    config.created_by_profile_field = 'name'
  } else {
    const strCol = cols.find(c => c.field_type === 'string' && c.column_name !== 'id')
    config.created_by_profile_field = strCol?.column_name || ''
  }
}

// 管理员角色选择（双向同步到 config.admin_role_ids）
const selectedAdminRoles = computed({
  get: () => {
    if (!config.admin_role_ids) return []
    return config.admin_role_ids.split(',').filter(s => s).map(s => parseInt(s.trim()))
  },
  set: (val: number[]) => {
    config.admin_role_ids = val.join(',')
  }
})

// 自动生成权限标识：父菜单权限标识（或path）+ 模块名
const autoPermission = computed(() => {
  if (!config.module_name) return ''
  
  // 查找父菜单
  const findMenu = (menus: Menu[], id: number): Menu | null => {
    for (const menu of menus) {
      if (menu.id === id) return menu
      if (menu.children) {
        const found = findMenu(menu.children, id)
        if (found) return found
      }
    }
    return null
  }
  
  if (menuConfig.parent_id && menuConfig.parent_id > 0) {
    const parentMenu = findMenu(menuTree.value, menuConfig.parent_id)
    if (parentMenu) {
      // 优先使用父菜单的权限标识，如果没有则使用path（去掉前导/）
      const parentPermission = parentMenu.permission || parentMenu.path?.replace(/^\//, '') || ''
      // 如果父菜单权限标识以:list结尾，去掉:list
      const prefix = parentPermission.replace(/:list$/, '')
      menuConfig.permission = prefix ? `${prefix}:${config.module_name}` : config.module_name
      return menuConfig.permission
    }
  }
  
  // 没有父菜单，直接使用模块名
  menuConfig.permission = config.module_name
  return config.module_name
})

// 字段表格列
const columnTableColumns = [
  { title: '', key: 'drag', width: 40, fixed: 'left' },
  { title: '字段名', key: 'column_name', width: 110 },
  { title: 'Go类型', key: 'field_type', width: 100 },
  { title: 'DB类型', key: 'db_type', width: 100 },
  { title: '长度', key: 'db_length', width: 70 },
  { title: '注释', key: 'comment', width: 110 },
  { title: '关联表', key: 'related_table', width: 200 },
  { title: '必填', key: 'is_required', width: 50 },
  { title: '搜索', key: 'is_searchable', width: 50 },
  { title: '搜索类型', key: 'search_type', width: 80 },
  { title: '列表', key: 'is_list_visible', width: 50 },
  { title: '排序', key: 'is_sortable', width: 50 },
  { title: '排序方式', key: 'sort_order', width: 80 },
  { title: '唯一', key: 'is_unique', width: 50 },
  { title: '表单', key: 'is_form_visible', width: 50 },
  { title: '组件', key: 'form_type', width: 100 },
  { title: '选项', key: 'options', width: 150 },
  { title: '操作', key: 'action', width: 100, fixed: 'right' }
]

// 一对多/多对多关联配置列（belongsTo 已在字段表中配置）
const relationColumns = [
  { title: '关联类型', key: 'relation_type', width: 160 },
  { title: '关联表', key: 'related_table', width: 140 },
  { title: '模块名', key: 'related_module', width: 120 },
  { title: '外键字段', key: 'foreign_key', width: 120 },
  { title: '显示字段', key: 'display_field', width: 100 },
  { title: '注释', key: 'comment', width: 100 },
  { title: '必填', key: 'is_required', width: 50 },
  { title: '中间表', key: 'join_table', width: 120 },
  { title: '轻量接口', key: 'use_options_api', width: 100 },
  { title: '操作', key: 'action', width: 60 }
]

const selectOptionColumns = [
  { title: '显示文本', key: 'label', width: 150 },
  { title: '值', key: 'value', width: 150 },
  { title: '操作', key: 'action', width: 80 }
]

// 工具函数
const toCamelCase = (str: string) => str.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase())
const toPascalCase = (str: string) => {
  const camel = toCamelCase(str)
  return camel.charAt(0).toUpperCase() + camel.slice(1)
}
const goTypeToTsType = (goType: string): string => {
  const map: Record<string, string> = {
    'string': 'string', 'int': 'number', 'int64': 'number', 'uint': 'number',
    'float64': 'number', 'bool': 'boolean', 'time.Time': 'string'
  }
  return map[goType] || 'string'
}

// 判断模块是否已生成
const isModuleGenerated = (moduleName: string) => generatedModules.value.includes(moduleName)

// 获取字段数量
const getColumnsCount = (record: SavedConfig) => {
  try {
    const parsed = JSON.parse(record.config_json)
    return parsed.columns?.length || 0
  } catch {
    return 0
  }
}

// 字段操作
const onColumnNameChange = (record: ColumnConfig) => {
  if (record.column_name) {
    record.field_name = toPascalCase(record.column_name)
    record.json_name = record.column_name
  }
}
const onFieldTypeChange = (record: ColumnConfig) => {
  record.ts_type = goTypeToTsType(record.field_type)
  // int类型自动选择数字组件
  if (['int', 'int64', 'uint', 'float64'].includes(record.field_type) && record.form_type === 'input') {
    record.form_type = 'number'
  }
}
// 表单类型改变时自动设置Go类型
const onFormTypeChange = (record: ColumnConfig) => {
  // 单个图片/文件用uint存储文件ID
  if (['image', 'file', 'upload'].includes(record.form_type)) {
    if (record.field_type === 'string') {
      record.field_type = 'uint'
      record.ts_type = 'number'
      message.info('单个图片/文件字段已自动设置为uint类型（存储文件ID）')
    }
  }
  // 多个图片/文件用string存储文件ID列表
  if (['images', 'files'].includes(record.form_type)) {
    if (record.field_type !== 'string') {
      record.field_type = 'string'
      record.ts_type = 'string'
      message.info('多个图片/文件字段已自动设置为string类型（存储文件ID列表）')
    }
  }
  // 日期/日期时间用time.Time
  if (['date', 'datetime'].includes(record.form_type)) {
    if (record.field_type !== 'time.Time') {
      record.field_type = 'time.Time'
      record.ts_type = 'string'
      message.info('日期字段已自动设置为time.Time类型')
    }
  }
  // 开关默认用int
  if (record.form_type === 'switch') {
    if (!['int', 'int8', 'int64', 'string', 'bool'].includes(record.field_type)) {
      record.field_type = 'int'
      record.ts_type = 'number'
    }
  }
  // 数字输入框默认用int
  if (record.form_type === 'number') {
    if (!['int', 'int8', 'int64', 'uint', 'float64'].includes(record.field_type)) {
      record.field_type = 'int'
      record.ts_type = 'number'
    }
  }
}
const addColumn = () => config.columns.push(createEmptyColumn())
const removeColumn = (index: number) => config.columns.splice(index, 1)

// 在指定位置前插入字段
const insertColumnBefore = (index: number) => {
  config.columns.splice(index, 0, createEmptyColumn())
}

// 在指定位置后插入字段
const insertColumnAfter = (index: number) => {
  config.columns.splice(index + 1, 0, createEmptyColumn())
}

// 拖拽排序状态
const dragState = reactive({
  draggingIndex: -1,
  dropIndex: -1,
  isDragging: false
})

// 拖拽手柄鼠标按下，允许拖拽
const handleDragHandleMouseDown = (index: number) => {
  dragState.draggingIndex = index
}

// 拖拽手柄鼠标松开
const handleDragHandleMouseUp = () => {
  dragState.draggingIndex = -1
}

// 自定义行属性
const customRow = (_record: any, index: number | undefined) => {
  return {
    draggable: dragState.draggingIndex === index,
    class: {
      'drag-over': dragState.dropIndex === index && dragState.isDragging,
      'dragging': dragState.draggingIndex === index && dragState.isDragging
    },
    onDragstart: (e: DragEvent) => {
      if (dragState.draggingIndex !== index) {
        e.preventDefault()
        return
      }
      dragState.isDragging = true
      e.dataTransfer!.effectAllowed = 'move'
      e.dataTransfer!.setData('text/plain', String(index))
    },
    onDragover: (e: DragEvent) => {
      e.preventDefault()
      e.dataTransfer!.dropEffect = 'move'
      if (index !== undefined) {
        dragState.dropIndex = index
      }
    },
    onDragleave: () => {
      dragState.dropIndex = -1
    },
    onDrop: (e: DragEvent) => {
      e.preventDefault()
      const fromIndex = parseInt(e.dataTransfer!.getData('text/plain'))
      const toIndex = index
      if (fromIndex !== toIndex && toIndex !== undefined) {
        const columns = [...config.columns]
        const [removed] = columns.splice(fromIndex, 1)
        columns.splice(toIndex, 0, removed)
        config.columns = columns
      }
      dragState.dropIndex = -1
      dragState.isDragging = false
      dragState.draggingIndex = -1
    },
    onDragend: () => {
      dragState.dropIndex = -1
      dragState.isDragging = false
      dragState.draggingIndex = -1
    }
  }
}

// 关联操作（hasMany/many2many）
const addRelation = () => {
  config.relations.push({
    relation_type: 'hasMany',
    related_table: '',
    related_module: '',
    related_model: '',
    foreign_key: '',
    reference_key: 'ID',
    join_table: '',
    display_field: 'name',
    comment: '',
    is_required: false,
    use_options_api: true,
    use_tree_layout: false
  })
}
const removeRelation = (index: number) => config.relations.splice(index, 1)

// 字段关联表选择变化时（belongsTo）
const onColumnRelatedTableChange = async (record: ColumnConfig) => {
  if (!record.related_table) {
    record.display_field = ''
    record.related_module = ''
    record.use_options_api = false
    record.use_tree_layout = false
    return
  }
  // 获取表注释作为默认comment
  const tableInfo = dbTables.value.find(t => t.table_name === record.related_table)
  if (tableInfo?.table_comment && !record.comment) {
    record.comment = tableInfo.table_comment
  }
  // 加载关联表字段
  if (!relationTableColumns.value[record.related_table]) {
    try {
      const res = await getTableColumns(record.related_table)
      relationTableColumns.value[record.related_table] = res.data || []
    } catch (e) {
      console.error('获取表字段失败', e)
    }
  }
  // 设置默认值
  record.display_field = 'name'
  record.use_options_api = true
  // 尝试自动设置display_field
  const cols = relationTableColumns.value[record.related_table] || []
  const nameCol = cols.find(c => c.column_name === 'name')
  if (!nameCol) {
    const strCol = cols.find(c => c.field_type === 'string' && c.column_name !== 'id')
    if (strCol) {
      record.display_field = strCol.column_name
    }
  }
}

// 字段关联配置弹窗
const columnRelationConfigVisible = ref(false)
const editingColumnRelationIndex = ref(-1)
const editingColumnRelation = reactive({
  related_module: '',
  display_field: 'name',
  use_options_api: true,
  use_tree_layout: false
})

const openColumnRelationConfig = (index: number) => {
  editingColumnRelationIndex.value = index
  const col = config.columns[index]
  Object.assign(editingColumnRelation, {
    related_module: col.related_module || '',
    display_field: col.display_field || 'name',
    use_options_api: col.use_options_api ?? true,
    use_tree_layout: col.use_tree_layout ?? false
  })
  columnRelationConfigVisible.value = true
}

const saveColumnRelationConfig = () => {
  if (editingColumnRelationIndex.value >= 0) {
    const col = config.columns[editingColumnRelationIndex.value]
    col.related_module = editingColumnRelation.related_module
    col.display_field = editingColumnRelation.display_field
    col.use_options_api = editingColumnRelation.use_options_api
    col.use_tree_layout = editingColumnRelation.use_tree_layout
    // 左树右表时自动开启轻量接口
    if (col.use_tree_layout) {
      col.use_options_api = true
    }
  }
  columnRelationConfigVisible.value = false
}

// 获取当前编辑字段的关联表字段列表
const editingColumnRelationColumns = computed(() => {
  if (editingColumnRelationIndex.value < 0) return []
  const col = config.columns[editingColumnRelationIndex.value]
  return relationTableColumns.value[col.related_table] || []
})

// 关联表选择变化时
const onRelatedTableChange = async (record: any) => {
  if (!record.related_table) {
    record.display_field = 'name'
    record.comment = ''
    return
  }
  // 设置关联模型名
  record.related_model = toPascalCase(record.related_table)
  // 获取表注释作为默认comment
  const tableInfo = dbTables.value.find(t => t.table_name === record.related_table)
  if (tableInfo?.table_comment && !record.comment) {
    record.comment = tableInfo.table_comment
  }
  // 加载关联表字段
  if (!relationTableColumns.value[record.related_table]) {
    try {
      const res = await getTableColumns(record.related_table)
      relationTableColumns.value[record.related_table] = res.data || []
    } catch (e) {
      console.error('获取表字段失败', e)
    }
  }
  // 尝试自动设置display_field（优先name, 其次第一个string字段）
  const cols = relationTableColumns.value[record.related_table] || []
  const nameCol = cols.find(c => c.column_name === 'name')
  if (nameCol) {
    record.display_field = 'name'
  } else {
    const strCol = cols.find(c => c.field_type === 'string' && c.column_name !== 'id')
    if (strCol) {
      record.display_field = strCol.column_name
    }
  }
}

// 获取关联表字段列表
const getRelationColumns = (tableName: string) => {
  return relationTableColumns.value[tableName] || []
}

// 显示字段选择变化时（不覆盖注释，注释用表注释更合适）
const onDisplayFieldChange = (_record: any) => {
  // 注释已经在选择关联表时用表注释填充，这里不再覆盖
}

// 预加载关联表字段（编辑/复制配置时）
const preloadRelationColumns = async (relations: any[]) => {
  if (!relations?.length) return
  const tables = relations.map(r => r.related_table).filter(t => t && !relationTableColumns.value[t])
  await Promise.all(tables.map(async (tableName) => {
    try {
      const res = await getTableColumns(tableName)
      relationTableColumns.value[tableName] = res.data || []
    } catch {
      // ignore
    }
  }))
}

// 下拉选项
const openSelectOptions = (index: number) => {
  editingColumnIndex.value = index
  const col = config.columns[index]
  editingSelectOptions.value = [...(col.select_options || [])]
  editingDictType.value = col.dict_type || ''
  useDictMode.value = !!col.dict_type
  selectOptionsVisible.value = true
}
const saveSelectOptions = () => {
  if (editingColumnIndex.value >= 0) {
    if (useDictMode.value && editingDictType.value) {
      // 使用字典模式
      config.columns[editingColumnIndex.value].dict_type = editingDictType.value
      config.columns[editingColumnIndex.value].select_options = []
    } else {
      // 使用手动选项模式
      config.columns[editingColumnIndex.value].dict_type = ''
      config.columns[editingColumnIndex.value].select_options = [...editingSelectOptions.value]
    }
  }
  selectOptionsVisible.value = false
}

// 开关值
const openSwitchValues = (index: number) => {
  editingColumnIndex.value = index
  const sv = config.columns[index].switch_values
  if (sv) {
    Object.assign(editingSwitchValues, sv)
  } else {
    Object.assign(editingSwitchValues, { active_value: '1', inactive_value: '0', active_text: '启用', inactive_text: '禁用' })
  }
  switchValuesVisible.value = true
}
const saveSwitchValues = () => {
  if (editingColumnIndex.value >= 0) {
    config.columns[editingColumnIndex.value].switch_values = { ...editingSwitchValues }
  }
  switchValuesVisible.value = false
}

// Options代码预览
const showOptionsCodePreview = (record: any) => {
  const tableName = record.related_table
  const modelName = toPascalCase(tableName)
  const routePath = tableName // 路由保持与表名一致（下划线格式）
  
  optionsCodeTable.value = tableName
  
  // Service 代码
  optionsCodeService.value = `// Get${modelName}Options 获取选项列表（带可选关联统计）
// excludeDeleted: 是否排除软删除数据（统计表有deleted_at字段时传true）
// countCreatedBy: 统计时按创建人过滤（数据隔离用，传当前用户ID，0表示不过滤）
func (s *${modelName}Service) Get${modelName}Options(displayField, countTable, countForeignKey string, excludeDeleted bool, countCreatedBy uint) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    if displayField == "" {
        displayField = "name"
    }
    
    // 无统计关联时，简单查询
    if countTable == "" || countForeignKey == "" {
        err := global.DB.Model(&model.${modelName}{}).
            Select("id, " + displayField + " as name").
            Order("id ASC").
            Find(&results).Error
        return results, err
    }
    
    // 有统计关联时，使用子查询
    subQuery := global.DB.Table(countTable).
        Select(countForeignKey + " as fk, COUNT(*) as cnt")
    
    // 排除软删除数据
    if excludeDeleted {
        subQuery = subQuery.Where("deleted_at IS NULL")
    }
    // 数据隔离：统计时按创建人过滤
    if countCreatedBy > 0 {
        subQuery = subQuery.Where("created_by = ?", countCreatedBy)
    }
    subQuery = subQuery.Group(countForeignKey)
    
    err := global.DB.Table("${tableName}").
        Select("${tableName}.id, ${tableName}." + displayField + " as name, COALESCE(sub.cnt, 0) as count").
        Joins("LEFT JOIN (?) as sub ON ${tableName}.id = sub.fk", subQuery).
        Order("${tableName}.id ASC").
        Find(&results).Error
    
    return results, err
}`

  // API 代码
  optionsCodeApi.value = `// Get${modelName}Options 获取选项列表
func (a *${modelName}Api) Get${modelName}Options(c *gin.Context) {
    displayField := c.DefaultQuery("display_field", "name")
    countTable := c.Query("count_table")
    countForeignKey := c.Query("count_field")
    excludeDeleted := c.Query("exclude_deleted") == "true"
    // 数据隔离：统计时按创建人过滤
    var countCreatedBy uint = 0
    if ccb := c.Query("count_created_by"); ccb != "" {
        if id, err := strconv.ParseUint(ccb, 10, 64); err == nil {
            countCreatedBy = uint(id)
        }
    }

    list, err := service.${modelName}.Get${modelName}Options(displayField, countTable, countForeignKey, excludeDeleted, countCreatedBy)
    if err != nil {
        response.Fail(c, "获取选项列表失败")
        return
    }
    response.OkWithData(c, list)
}`

  // Router 代码
  optionsCodeRouter.value = `// 在 RegisterPrivateRoutes 中添加
R(rg, "GET", "/${routePath}/options", m.Name(), "选项列表", v1.${modelName}.Get${modelName}Options, registry.WithAuth())`

  // 前端 API 代码
  optionsCodeFrontend.value = `// 获取选项列表
export function get${modelName}Options(params?: { display_field?: string; count_table?: string; count_field?: string; exclude_deleted?: boolean; count_created_by?: number }) {
  return request.get<any, ApiResponse<OptionItem[]>>('/${routePath}/options', { params })
}

// OptionItem 类型定义（如果没有）
export interface OptionItem {
  id: number
  name: string
  count?: number
}`

  optionsCodeVisible.value = true
}

// SQL导入
const openImportSQL = () => {
  importSQLText.value = ''
  importSQLVisible.value = true
}

// JSON配置导入
const openImportJSON = () => {
  importJSONText.value = ''
  importJSONVisible.value = true
}

// 解析JSON并导入
const handleImportJSON = async () => {
  const jsonStr = importJSONText.value.trim()
  if (!jsonStr) {
    message.warning('请输入JSON配置')
    return
  }
  try {
    const parsed = JSON.parse(jsonStr) as GeneratorConfig
    
    // 验证必要字段
    if (!parsed.table_name && !parsed.module_name) {
      message.warning('JSON配置缺少table_name或module_name字段')
      return
    }
    
    // 修正关联关系：左树右表必须开启轻量接口
    parsed.relations?.forEach(r => {
      if (r.use_tree_layout) r.use_options_api = true
    })
    
    // 填充基础配置
    Object.assign(config, {
      id: undefined, // 导入的配置作为新配置
      table_name: parsed.table_name || '',
      module_name: parsed.module_name || '',
      description: parsed.description || '',
      author: parsed.author || '',
      generate_backend: parsed.generate_backend !== false,
      generate_frontend: parsed.generate_frontend !== false,
      generate_sql: parsed.generate_sql !== false,
      frontend_path: parsed.frontend_path || '',
      has_created_at: parsed.has_created_at !== false,
      has_updated_at: parsed.has_updated_at !== false,
      has_deleted_at: parsed.has_deleted_at === true,
      has_created_by: parsed.has_created_by === true,
      data_isolation: parsed.data_isolation === true,
      admin_role_ids: parsed.admin_role_ids || '',
      has_audit: parsed.has_audit === true,
      generate_frontend_api: parsed.generate_frontend_api === true,
      link_to_user: parsed.link_to_user === true,
      profile_name: parsed.profile_name || '',
      profile_icon: parsed.profile_icon || '',
      profile_role_code: parsed.profile_role_code || '',
      enable_import_export: parsed.enable_import_export === true,
      columns: parsed.columns || [],
      relations: parsed.relations || [],
      menu_config: parsed.menu_config || null
    })
    
    // 填充菜单配置
    if (parsed.menu_config) {
      Object.assign(menuConfig, parsed.menu_config)
    } else {
      Object.assign(menuConfig, { parent_id: 0, menu_name: '', menu_icon: '', menu_sort: 0, permission: '' })
    }
    
    // 预加载关联表字段
    await preloadRelationColumns(parsed.relations || [])
    
    importJSONVisible.value = false
    message.success(`已导入配置：${parsed.module_name || parsed.table_name}（${parsed.columns?.length || 0}个字段，${parsed.relations?.length || 0}个关联）`)
  } catch (e: any) {
    message.error('JSON解析失败: ' + (e.message || '格式错误'))
  }
}

// 导出JSON配置
const handleExportJSON = async () => {
  const exportConfig = buildConfig()
  // 移除id字段（导出的配置不需要id）
  delete exportConfig.id
  const jsonStr = JSON.stringify(exportConfig, null, 2)
  try {
    await navigator.clipboard.writeText(jsonStr)
    message.success('配置JSON已复制到剪贴板')
  } catch {
    // 如果剪贴板写入失败，显示弹窗让用户手动复制
    importJSONText.value = jsonStr
    importJSONVisible.value = true
    message.info('请手动复制以下配置')
  }
}

// 解析SQL并导入
const handleImportSQL = () => {
  const sql = importSQLText.value.trim()
  if (!sql) {
    message.warning('请输入SQL语句')
    return
  }
  try {
    const result = parseCreateTableSQL(sql)
    // 填充表名和模块名
    if (result.tableName) {
      config.table_name = result.tableName
      config.module_name = result.tableName.replace(/^[a-z]+_/, '') // 去掉前缀
    }
    if (result.tableComment) {
      config.description = result.tableComment
    }
    // 替换字段配置
    config.columns = result.columns
    importSQLVisible.value = false
    message.success(`已导入 ${result.columns.length} 个字段`)
  } catch (e: any) {
    message.error('SQL解析失败: ' + (e.message || '格式错误'))
  }
}

// SQL解析函数
const parseCreateTableSQL = (sql: string): { tableName: string; tableComment: string; columns: ColumnConfig[] } => {
  // 提取表名
  const tableNameMatch = sql.match(/CREATE\s+TABLE\s+[`'"]?([\w]+)[`'"]?/i)
  const tableName = tableNameMatch ? tableNameMatch[1] : ''
  
  // 提取表注释
  const tableCommentMatch = sql.match(/COMMENT\s*=?\s*['"]([^'"]+)['"]\s*;?\s*$/im)
  const tableComment = tableCommentMatch ? tableCommentMatch[1] : ''
  
  // 提取字段定义部分 - 匹配第一个(和最后一个)之间的内容
  const columnsMatch = sql.match(/\(([\s\S]+)\)[^)]*$/)
  if (!columnsMatch) {
    throw new Error('无法解析字段定义')
  }
  
  const columnsSection = columnsMatch[1]
  const columns: ColumnConfig[] = []
  
  // 跳过的系统字段
  const skipColumns = ['id', 'created_at', 'updated_at', 'deleted_at', 'created_by']
  
  // 按行分割，处理每个字段
  const lines = columnsSection.split(/,(?=\s*(?:`|\w|PRIMARY|KEY|INDEX|UNIQUE|CONSTRAINT|FOREIGN))/i)
  
  for (const line of lines) {
    const trimmed = line.trim()
    // 跳过索引、主键、约束等
    if (/^(PRIMARY|KEY|INDEX|UNIQUE|CONSTRAINT|FOREIGN)/i.test(trimmed)) {
      continue
    }
    
    // 解析字段: `column_name` TYPE(...) ... COMMENT 'xxx'
    const fieldMatch = trimmed.match(/^[`'"]?([\w]+)[`'"]?\s+([\w]+)(?:\s*\(([\d,]+)\))?([^]*?)$/i)
    if (!fieldMatch) continue
    
    const columnName = fieldMatch[1]
    const dbType = fieldMatch[2].toUpperCase()
    const dbLength = fieldMatch[3] ? parseInt(fieldMatch[3].split(',')[0]) : 0
    const rest = fieldMatch[4] || ''
    
    // 跳过系统字段
    if (skipColumns.includes(columnName.toLowerCase())) {
      continue
    }
    
    // 提取注释
    const commentMatch = rest.match(/COMMENT\s+['"]([^'"]*)['"]\s*$/i)
    const comment = commentMatch ? commentMatch[1] : ''
    
    // 检测NOT NULL
    const isRequired = /NOT\s+NULL/i.test(rest) && !/DEFAULT/i.test(rest)
    
    // DB类型转Go类型
    const { goType, tsType, formType } = dbTypeToGoType(dbType, dbLength, columnName)
    
    columns.push({
      column_name: columnName,
      field_name: toPascalCase(columnName),
      field_type: goType,
      json_name: columnName,
      ts_type: tsType,
      comment: comment,
      db_type: dbType,
      db_length: dbLength,
      default_value: '',
      is_primary_key: false,
      is_required: isRequired,
      is_searchable: false,
      search_type: 'eq',
      is_list_visible: true,
      is_form_visible: true,
      is_sortable: false,
      sort_order: 'asc',
      is_unique: false,
      form_type: formType,
      dict_type: '',
      select_options: [],
      switch_values: null
    })
  }
  
  return { tableName, tableComment, columns }
}

// DB类型转Go类型
const dbTypeToGoType = (dbType: string, _length: number, columnName: string): { goType: string; tsType: string; formType: string } => {
  const type = dbType.toUpperCase()
  
  // 根据字段名推断组件类型
  let formType = 'input'
  if (/status|type|state|level|difficulty/i.test(columnName)) {
    formType = 'select'
  } else if (/content|description|remark|analysis/i.test(columnName)) {
    formType = 'textarea'
  } else if (/_time$/i.test(columnName)) {
    formType = 'datetime'
  } else if (/_date$/i.test(columnName)) {
    formType = 'date'
  }
  
  switch (type) {
    case 'TINYINT':
    case 'SMALLINT':
    case 'MEDIUMINT':
    case 'INT':
    case 'INTEGER':
      return { goType: 'int', tsType: 'number', formType: formType === 'input' ? 'number' : formType }
    case 'BIGINT':
      // 检查是否是无符号(ID字段)
      if (/_id$/i.test(columnName)) {
        return { goType: 'uint', tsType: 'number', formType: 'number' }
      }
      return { goType: 'int64', tsType: 'number', formType: formType === 'input' ? 'number' : formType }
    case 'FLOAT':
    case 'DOUBLE':
    case 'DECIMAL':
      return { goType: 'float64', tsType: 'number', formType: 'number' }
    case 'BOOL':
    case 'BOOLEAN':
      return { goType: 'bool', tsType: 'boolean', formType: 'switch' }
    case 'DATE':
      return { goType: 'time.Time', tsType: 'string', formType: 'date' }
    case 'DATETIME':
    case 'TIMESTAMP':
      return { goType: 'time.Time', tsType: 'string', formType: 'datetime' }
    case 'TEXT':
    case 'MEDIUMTEXT':
    case 'LONGTEXT':
      return { goType: 'string', tsType: 'string', formType: 'textarea' }
    default:
      return { goType: 'string', tsType: 'string', formType }
  }
}

// 重置配置
const resetConfig = () => {
  Object.assign(config, {
    id: undefined, table_name: '', module_name: '', description: '', author: '',
    generate_backend: true, generate_frontend: true, generate_sql: true,
    frontend_path: '', has_created_at: true, has_updated_at: true, has_deleted_at: false, has_created_by: false,
    created_by_profile_table: '', created_by_profile_field: '',
    data_isolation: false, admin_role_ids: '', has_audit: false, generate_frontend_api: false,
    link_to_user: false, profile_name: '', profile_icon: '', profile_role_code: '',
    enable_import_export: false,
    columns: [], relations: [], menu_config: null, stats_config: null
  })
  Object.assign(menuConfig, { parent_id: 0, menu_name: '', menu_icon: '', menu_sort: 0, permission: '' })
  Object.assign(statsConfig, { enabled: false, charts: [], time_field: '' })
}

// 新增配置
const handleAdd = () => {
  resetConfig()
  drawerTitle.value = '新增配置'
  activeTab.value = '1'
  drawerVisible.value = true
}

// 编辑配置
const handleEdit = async (record: SavedConfig) => {
  try {
    const parsed = JSON.parse(record.config_json) as GeneratorConfig
    // 修正关联关系：左树右表必须开启轻量接口
    parsed.relations?.forEach(r => {
      if (r.use_tree_layout) r.use_options_api = true
    })
    // 保存配置ID用于更新
    parsed.id = record.id
    Object.assign(config, parsed)
    if (parsed.menu_config) {
      Object.assign(menuConfig, parsed.menu_config)
    } else {
      Object.assign(menuConfig, { parent_id: 0, menu_name: '', menu_icon: '', menu_sort: 0, permission: '' })
    }
    // 加载统计配置
    if (parsed.stats_config) {
      Object.assign(statsConfig, parsed.stats_config)
    } else {
      Object.assign(statsConfig, { enabled: false, charts: [], time_field: '' })
    }
    // 预加载关联表字段
    await preloadRelationColumns(parsed.relations)
    drawerTitle.value = `编辑配置 - ${record.module_name}`
    activeTab.value = '1'
    drawerVisible.value = true
  } catch {
    message.error('配置解析失败')
  }
}

// 复制配置
const handleCopy = async (record: SavedConfig) => {
  try {
    const parsed = JSON.parse(record.config_json) as GeneratorConfig
    // 修正关联关系：左树右表必须开启轻量接口
    parsed.relations?.forEach(r => {
      if (r.use_tree_layout) r.use_options_api = true
    })
    // 复制时清空模块名和表名，让用户填写新的
    Object.assign(config, {
      ...parsed,
      id: undefined,  // 复制时清空ID，作为新配置保存
      table_name: '',
      module_name: '',
      description: parsed.description + '(复制)'
    })
    if (parsed.menu_config) {
      Object.assign(menuConfig, {
        ...parsed.menu_config,
        menu_name: ''
      })
    } else {
      Object.assign(menuConfig, { parent_id: 0, menu_name: '', menu_icon: '', menu_sort: 0, permission: '' })
    }
    // 加载统计配置
    if (parsed.stats_config) {
      Object.assign(statsConfig, parsed.stats_config)
    } else {
      Object.assign(statsConfig, { enabled: false, charts: [], time_field: '' })
    }
    // 预加载关联表字段
    await preloadRelationColumns(parsed.relations)
    drawerTitle.value = '新增配置(复制)'
    activeTab.value = '1'
    drawerVisible.value = true
    message.success('已复制配置，请修改表名和模块名称')
  } catch {
    message.error('配置解析失败')
  }
}

// 构建配置
const buildConfig = (): GeneratorConfig => {
  // 修正关联关系：左树右表必须开启轻量接口
  config.relations?.forEach(r => {
    if (r.use_tree_layout) r.use_options_api = true
  })
  return {
    ...config,
    menu_config: menuConfig.menu_name ? menuConfig : null,
    stats_config: statsConfig.enabled ? { ...statsConfig } : null
  }
}

// 保存配置
const handleSaveConfig = async () => {
  if (!config.table_name || !config.module_name) {
    message.warning('请填写表名和模块名称')
    return
  }
  saveLoading.value = true
  try {
    await saveConfigApi(buildConfig())
    message.success('配置保存成功')
    fetchSavedConfigs()
  } finally {
    saveLoading.value = false
  }
}

// 预览代码
const handlePreview = async () => {
  previewLoading.value = true
  try {
    const res = await previewCode(buildConfig())
    previewFiles.value = res.data.files
    if (previewFiles.value.length > 0) {
      previewTab.value = previewFiles.value[0].path
    }
    previewVisible.value = true
  } finally {
    previewLoading.value = false
  }
}

// 导出配置JSON（从列表）
const handleExportConfigJSON = async (record: SavedConfig) => {
  try {
    const parsed = JSON.parse(record.config_json) as GeneratorConfig
    // 移除id字段
    delete parsed.id
    const jsonStr = JSON.stringify(parsed, null, 2)
    await navigator.clipboard.writeText(jsonStr)
    message.success(`配置 "${record.module_name}" 已复制到剪贴板`)
  } catch {
    message.error('导出失败')
  }
}

// 显示E-R图
const handleShowERDiagram = (record: SavedConfig) => {
  try {
    const parsed = JSON.parse(record.config_json) as GeneratorConfig
    
    // 构建实体列表
    const entity: EREntity = {
      name: parsed.table_name,
      comment: parsed.description || parsed.module_name,
      columns: parsed.columns.map(col => ({
        name: col.column_name,
        comment: col.comment || col.column_name,
        isPrimary: col.is_primary_key || col.column_name === 'id'
      }))
    }
    
    // 收集所有实体（包含关联表）
    const entityMap = new Map<string, EREntity>()
    entityMap.set(parsed.table_name, entity)
    
    // 构建关联列表
    const relationList: ERRelation[] = []
    
    if (parsed.relations && parsed.relations.length > 0) {
      parsed.relations.forEach(rel => {
        // 添加关联表作为实体
        if (rel.related_table && !entityMap.has(rel.related_table)) {
          // 尝试从已保存的配置中获取关联表的完整字段
          const relatedConfig = savedConfigs.value.find(c => c.table_name === rel.related_table)
          let relatedColumns: { name: string; comment: string; isPrimary?: boolean }[] = []
          
          if (relatedConfig) {
            // 从已保存配置中获取完整字段
            try {
              const relatedParsed = JSON.parse(relatedConfig.config_json) as GeneratorConfig
              relatedColumns = relatedParsed.columns.map(col => ({
                name: col.column_name,
                comment: col.comment || col.column_name,
                isPrimary: col.is_primary_key || col.column_name === 'id'
              }))
            } catch {
              // 解析失败，使用默认字段
            }
          }
          
          // 如果没有找到配置，使用默认字段
          if (relatedColumns.length === 0) {
            relatedColumns = [
              { name: 'id', comment: 'ID', isPrimary: true },
              { name: rel.display_field || 'name', comment: rel.display_field || '名称', isPrimary: false }
            ]
          }
          
          entityMap.set(rel.related_table, {
            name: rel.related_table,
            comment: rel.comment || relatedConfig?.description || rel.related_model || rel.related_table,
            columns: relatedColumns
          })
        }
        
        // 根据关联类型确定基数和关系名称
        let fromCard = '1'
        let toCard = 'n'
        let relationName = '关联'
        
        switch (rel.relation_type) {
          case 'belongsTo':
            // 属于：本表多对一
            fromCard = 'n'
            toCard = '1'
            relationName = '属于'
            break
          case 'hasMany':
            // 一对多：本表一对多
            fromCard = '1'
            toCard = 'n'
            relationName = '包含'
            break
          case 'many2many':
            // 多对多
            fromCard = 'm'
            toCard = 'n'
            relationName = '关联'
            break
        }
        
        // 获取关联表的显示名称
        const relatedEntityName = rel.comment || rel.related_model || rel.related_table
        
        relationList.push({
          name: relationName,
          from: parsed.description || parsed.module_name,
          to: relatedEntityName,
          fromCardinality: fromCard,
          toCardinality: toCard
        })
      })
    }
    
    erEntities.value = Array.from(entityMap.values())
    erRelations.value = relationList
    erDiagramVisible.value = true
  } catch (e) {
    console.error('E-R图生成失败:', e)
    message.error('配置解析失败')
  }
}

// 从配置预览
const handlePreviewFromConfig = async (record: SavedConfig) => {
  try {
    const parsed = JSON.parse(record.config_json) as GeneratorConfig
    // 修正关联关系：左树右表必须开启轻量接口
    parsed.relations?.forEach(r => {
      if (r.use_tree_layout) r.use_options_api = true
    })
    previewLoading.value = true
    const res = await previewCode(parsed)
    previewFiles.value = res.data.files
    if (previewFiles.value.length > 0) {
      previewTab.value = previewFiles.value[0].path
    }
    previewVisible.value = true
  } catch {
    message.error('配置解析失败')
  } finally {
    previewLoading.value = false
  }
}

// 生成代码
const handleGenerate = async () => {
  if (!config.table_name || !config.module_name) {
    message.warning('请填写表名和模块名称')
    return
  }
  // 生成前端代码时才需要菜单配置
  if (config.generate_frontend && !menuConfig.menu_name) {
    message.warning('请填写菜单配置')
    activeTab.value = '4'
    return
  }
  if (generatedModules.value.includes(config.module_name)) {
    message.error(`模块 "${config.module_name}" 已存在，请先删除后再生成！`)
    return
  }
  generateLoading.value = true
  try {
    // 先保存配置
    await saveConfigApi(buildConfig())
    // 再生成代码
    await generateCode(buildConfig())
    message.success('代码生成成功！请重启后端服务以加载新模块。')
    fetchSavedConfigs()
    fetchModules()
    drawerVisible.value = false
  } catch (e: any) {
    message.error(e?.response?.data?.msg || e?.message || '生成失败')
  } finally {
    generateLoading.value = false
  }
}

// 从配置生成
const handleGenerateFromConfig = async (record: SavedConfig) => {
  if (generatedModules.value.includes(record.module_name)) {
    message.error(`模块 "${record.module_name}" 已存在，请先删除后再生成！`)
    return
  }
  try {
    const parsed = JSON.parse(record.config_json) as GeneratorConfig
    // 修正关联关系：左树右表必须开启轻量接口
    parsed.relations?.forEach(r => {
      if (r.use_tree_layout) r.use_options_api = true
    })
    // 生成前端代码时才需要菜单配置
    if (parsed.generate_frontend && !parsed.menu_config?.menu_name) {
      message.warning('请先编辑配置并填写菜单配置')
      return
    }
    generateLoading.value = true
    await generateCode(parsed)
    message.success('代码生成成功！请重启后端服务以加载新模块。')
    fetchModules()
  } catch (e: any) {
    message.error(e?.response?.data?.msg || e?.message || '生成失败')
  } finally {
    generateLoading.value = false
  }
}

// 执行SQL
const handleExecuteSQL = async (sql: string) => {
  executeSqlLoading.value = true
  try {
    await executeSQL(sql)
    message.success('SQL执行成功！表已创建。')
  } catch (e: any) {
    message.error(e?.response?.data?.msg || e?.message || 'SQL执行失败')
  } finally {
    executeSqlLoading.value = false
  }
}

// 获取显示路径
const getDisplayPath = (fullPath: string) => {
  const parts = fullPath.replace(/\\/g, '/').split('/')
  const idx = parts.lastIndexOf('src')
  if (idx >= 0) return parts.slice(idx).join('/')
  const sIdx = parts.lastIndexOf('go-base-server')
  if (sIdx >= 0) return parts.slice(sIdx + 1).join('/')
  return parts.slice(-2).join('/')
}

// 获取代码语言类名（用于高亮样式）
const getCodeLanguage = (filePath: string) => {
  const ext = filePath.split('.').pop()?.toLowerCase()
  const langMap: Record<string, string> = {
    'go': 'language-go',
    'ts': 'language-typescript',
    'vue': 'language-vue',
    'sql': 'language-sql'
  }
  return langMap[ext || ''] || 'language-plaintext'
}

// 复制代码
const handleCopyCode = async (content: string) => {
  try {
    await navigator.clipboard.writeText(content)
    message.success('代码已复制到剪贴板')
  } catch {
    message.error('复制失败')
  }
}

// 检测字段排序是否变化
const checkSortChanged = (originalCols: ColumnConfig[], currentCols: ColumnConfig[]) => {
  // 获取共同存在的字段（排除空字段）
  const commonFields = currentCols
    .filter(c => c.column_name && originalCols.find(o => o.column_name === c.column_name))
    .map(c => c.column_name)
  
  const originalOrder = originalCols
    .filter(c => commonFields.includes(c.column_name))
    .map(c => c.column_name)
  
  // 比较顺序
  for (let i = 0; i < commonFields.length; i++) {
    if (commonFields[i] !== originalOrder[i]) {
      return true
    }
  }
  return false
}

// 显示变更指南
const handleShowChangeGuide = () => {
  // 获取已保存的原始配置
  const savedConfig = savedConfigs.value.find(c => c.module_name === config.module_name)
  if (!savedConfig) {
    message.warning('未找到已保存的配置')
    return
  }

  let originalConfig: GeneratorConfig
  try {
    originalConfig = JSON.parse(savedConfig.config_json)
  } catch {
    message.error('解析原始配置失败')
    return
  }

  const categories: ChangeCategory[] = []
  const modelName = toPascalCase(config.table_name)
  const moduleName = config.module_name

  // 对比字段变化
  const originalCols = originalConfig.columns || []
  const currentCols = config.columns || []
  
  // 新增的字段（排除空字段）
  const addedCols = currentCols.filter(c => c.column_name && !originalCols.find(o => o.column_name === c.column_name))
  // 检测排序变化
  const sortChanged = checkSortChanged(originalCols, currentCols)
  // 搜索状态变化的字段
  const searchChanged = currentCols.filter(c => {
    const orig = originalCols.find(o => o.column_name === c.column_name)
    return orig && orig.is_searchable !== c.is_searchable
  })
  // 列表显示状态变化的字段
  const listVisibleChanged = currentCols.filter(c => {
    const orig = originalCols.find(o => o.column_name === c.column_name)
    return orig && orig.is_list_visible !== c.is_list_visible
  })
  // 表单显示状态变化的字段
  const formVisibleChanged = currentCols.filter(c => {
    const orig = originalCols.find(o => o.column_name === c.column_name)
    return orig && orig.is_form_visible !== c.is_form_visible
  })

  const newSearchCols = searchChanged.filter(c => c.is_searchable)
  const newListCols = listVisibleChanged.filter(c => c.is_list_visible)
  const newFormCols = formVisibleChanged.filter(c => c.is_form_visible)

  // === 变更概览 Tab ===
  const summaryItems: string[] = []
  if (addedCols.length > 0) {
    summaryItems.push(`新增 ${addedCols.length} 个字段：${addedCols.map(c => c.comment || c.column_name).join('、')}`)
  }
  if (newSearchCols.length > 0) {
    summaryItems.push(`${newSearchCols.length} 个字段开启搜索：${newSearchCols.map(c => c.comment || c.column_name).join('、')}`)
  }
  if (newListCols.length > 0) {
    summaryItems.push(`${newListCols.length} 个字段开启列表显示：${newListCols.map(c => c.comment || c.column_name).join('、')}`)
  }
  if (newFormCols.length > 0) {
    summaryItems.push(`${newFormCols.length} 个字段开启表单显示：${newFormCols.map(c => c.comment || c.column_name).join('、')}`)
  }
  if (sortChanged) {
    summaryItems.push(`字段排序已变更`)
  }

  if (summaryItems.length === 0) {
    message.info('配置无变化')
    return
  }

  categories.push({
    key: 'summary',
    label: `📝 变更概览`,
    summary: summaryItems,
    guides: []
  })

  // === 新增字段 Tab ===
  if (addedCols.length > 0) {
    const guides: ChangeGuide[] = []
    const formCols = addedCols.filter(c => c.is_form_visible)
    const listCols = addedCols.filter(c => c.is_list_visible)

    // Model
    guides.push({
      title: `模型新增字段`,
      description: `在 ${modelName} 结构体中添加`,
      file: `model/${moduleName}.go`,
      code: addedCols.map(c => {
        const fieldName = toPascalCase(c.column_name)
        return `\t${fieldName} ${c.field_type} \`json:"${c.column_name}" gorm:"comment:${c.comment}"\``
      }).join('\n')
    })

    if (formCols.length > 0) {
      // Request
      guides.push({
        title: `请求结构体新增字段`,
        description: `在 Create${modelName}Request 和 Update${modelName}Request 中添加`,
        file: `model/request/${moduleName}_request.go`,
        code: formCols.map(c => {
          const fieldName = toPascalCase(c.column_name)
          const reqType = c.is_required ? c.field_type : `*${c.field_type}`
          const bindTag = c.is_required ? `binding:"required"` : ''
          return `\t${fieldName} ${reqType} \`json:"${c.column_name}" ${bindTag}\` // ${c.comment}`
        }).join('\n')
      })

      // Service Create
      guides.push({
        title: `Service Create 新增字段`,
        description: `在 Create${modelName} 函数的 data := model.${modelName}{...} 中添加`,
        file: `service/${moduleName}.go`,
        code: formCols.map(c => {
          const fieldName = toPascalCase(c.column_name)
          const deref = c.is_required && ['int', 'int64', 'uint', 'float64'].includes(c.field_type) ? '*' : ''
          return `\t\t${fieldName}: ${deref}req.${fieldName},`
        }).join('\n')
      })

      // Service Update
      guides.push({
        title: `Service Update 新增字段`,
        description: `在 Update${modelName} 函数的 updates := map[string]interface{}{...} 中添加`,
        file: `service/${moduleName}.go`,
        code: formCols.map(c => {
          const fieldName = toPascalCase(c.column_name)
          return `\t\t"${c.column_name}": req.${fieldName},`
        }).join('\n')
      })
    }

    // 前端类型
    guides.push({
      title: `前端类型新增字段`,
      description: `在 ${modelName} 和 Create${modelName}Request 接口中添加`,
      file: `src/api/types/${moduleName}.ts`,
      code: addedCols.map(c => `  ${c.column_name}${c.is_required ? '' : '?'}: ${c.ts_type}`).join('\n')
    })

    // 前端表格列
    if (listCols.length > 0) {
      guides.push({
        title: `前端新增表格列`,
        description: `在 columns 数组中添加`,
        file: `src/views/xxx/${moduleName}/index.vue`,
        code: listCols.map(c => `{ title: '${c.comment}', dataIndex: '${c.column_name}', key: '${c.column_name}' }`).join(',\n')
      })
    }

    // 前端表单项
    if (formCols.length > 0) {
      guides.push({
        title: `前端新增表单项`,
        description: `在表单中添加`,
        file: `src/views/xxx/${moduleName}/Form.vue`,
        code: formCols.map(c => {
          const rules = c.is_required ? `:rules="[{ required: true, message: '请输入${c.comment}' }]"` : ''
          if (c.form_type === 'textarea') return `<a-form-item label="${c.comment}" name="${c.column_name}" ${rules}>\n  <a-textarea v-model:value="formData.${c.column_name}" placeholder="请输入${c.comment}" />\n</a-form-item>`
          if (c.form_type === 'number') return `<a-form-item label="${c.comment}" name="${c.column_name}" ${rules}>\n  <a-input-number v-model:value="formData.${c.column_name}" style="width: 100%" />\n</a-form-item>`
          if (c.form_type === 'select') return `<a-form-item label="${c.comment}" name="${c.column_name}" ${rules}>\n  <a-select v-model:value="formData.${c.column_name}" placeholder="请选择" />\n</a-form-item>`
          if (c.form_type === 'switch') return `<a-form-item label="${c.comment}" name="${c.column_name}">\n  <a-switch v-model:checked="formData.${c.column_name}" />\n</a-form-item>`
          return `<a-form-item label="${c.comment}" name="${c.column_name}" ${rules}>\n  <a-input v-model:value="formData.${c.column_name}" placeholder="请输入${c.comment}" />\n</a-form-item>`
        }).join('\n\n')
      })
    }

    categories.push({
      key: 'added',
      label: `➕ 新增字段 (${addedCols.length})`,
      summary: [`新增字段：${addedCols.map(c => c.comment || c.column_name).join('、')}`],
      guides
    })
  }

  // === 搜索变化 Tab ===
  if (newSearchCols.length > 0) {
    const guides: ChangeGuide[] = []

    guides.push({
      title: `Service 新增搜索条件`,
      description: `在 Get${modelName}List 函数的查询条件中添加`,
      file: `service/${moduleName}.go`,
      code: newSearchCols.map(c => {
        const fieldName = toPascalCase(c.column_name)
        if (c.search_type === 'like' || c.field_type === 'string') {
          return `if req.${fieldName} != "" {\n\tdb = db.Where("${c.column_name} LIKE ?", "%"+req.${fieldName}+"%")\n}`
        }
        return `if req.${fieldName} != nil {\n\tdb = db.Where("${c.column_name} = ?", *req.${fieldName})\n}`
      }).join('\n\n')
    })

    guides.push({
      title: `请求结构体新增搜索字段`,
      description: `在 ${modelName}ListRequest 中添加`,
      file: `model/request/${moduleName}_request.go`,
      code: newSearchCols.map(c => {
        const fieldName = toPascalCase(c.column_name)
        if (c.field_type === 'string') return `\t${fieldName} string \`form:"${c.column_name}"\` // ${c.comment}`
        return `\t${fieldName} *${c.field_type} \`form:"${c.column_name}"\` // ${c.comment}`
      }).join('\n')
    })

    guides.push({
      title: `前端 Query 类型新增字段`,
      description: `在 ${modelName}Query 接口中添加`,
      file: `src/api/types/${moduleName}.ts`,
      code: newSearchCols.map(c => `  ${c.column_name}?: ${c.ts_type}`).join('\n')
    })

    guides.push({
      title: `前端新增搜索表单项`,
      description: `在搜索表单中添加`,
      file: `src/views/xxx/${moduleName}/index.vue`,
      code: newSearchCols.map(c => {
        if (c.form_type === 'select') return `<a-form-item label="${c.comment}">\n  <a-select v-model:value="searchForm.${c.column_name}" placeholder="请选择" allowClear style="width: 160px" />\n</a-form-item>`
        return `<a-form-item label="${c.comment}">\n  <a-input v-model:value="searchForm.${c.column_name}" placeholder="请输入" allowClear style="width: 200px" />\n</a-form-item>`
      }).join('\n\n')
    })

    categories.push({
      key: 'search',
      label: `🔍 搜索变化 (${newSearchCols.length})`,
      summary: [`开启搜索：${newSearchCols.map(c => c.comment || c.column_name).join('、')}`],
      guides
    })
  }

  // === 列表显示变化 Tab ===
  if (newListCols.length > 0) {
    categories.push({
      key: 'list',
      label: `📊 列表变化 (${newListCols.length})`,
      summary: [`开启列表显示：${newListCols.map(c => c.comment || c.column_name).join('、')}`],
      guides: [{
        title: `前端新增表格列`,
        description: `在 columns 数组中添加`,
        file: `src/views/xxx/${moduleName}/index.vue`,
        code: newListCols.map(c => `{ title: '${c.comment}', dataIndex: '${c.column_name}', key: '${c.column_name}' }`).join(',\n')
      }]
    })
  }

  // === 表单显示变化 Tab ===
  if (newFormCols.length > 0) {
    const guides: ChangeGuide[] = []

    guides.push({
      title: `请求结构体新增字段`,
      description: `在 Create${modelName}Request 和 Update${modelName}Request 中添加`,
      file: `model/request/${moduleName}_request.go`,
      code: newFormCols.map(c => {
        const fieldName = toPascalCase(c.column_name)
        const reqType = c.is_required ? c.field_type : `*${c.field_type}`
        const bindTag = c.is_required ? `binding:"required"` : ''
        return `\t${fieldName} ${reqType} \`json:"${c.column_name}" ${bindTag}\` // ${c.comment}`
      }).join('\n')
    })

    guides.push({
      title: `Service Create 新增字段`,
      description: `在 Create${modelName} 函数的 data := model.${modelName}{...} 中添加`,
      file: `service/${moduleName}.go`,
      code: newFormCols.map(c => {
        const fieldName = toPascalCase(c.column_name)
        const deref = c.is_required && ['int', 'int64', 'uint', 'float64'].includes(c.field_type) ? '*' : ''
        return `\t\t${fieldName}: ${deref}req.${fieldName},`
      }).join('\n')
    })

    guides.push({
      title: `Service Update 新增字段`,
      description: `在 Update${modelName} 函数的 updates := map[string]interface{}{...} 中添加`,
      file: `service/${moduleName}.go`,
      code: newFormCols.map(c => {
        const fieldName = toPascalCase(c.column_name)
        return `\t\t"${c.column_name}": req.${fieldName},`
      }).join('\n')
    })

    guides.push({
      title: `前端类型新增字段`,
      description: `在 Create${modelName}Request 接口中添加`,
      file: `src/api/types/${moduleName}.ts`,
      code: newFormCols.map(c => `  ${c.column_name}${c.is_required ? '' : '?'}: ${c.ts_type}`).join('\n')
    })

    guides.push({
      title: `前端新增表单项`,
      description: `在表单中添加`,
      file: `src/views/xxx/${moduleName}/Form.vue`,
      code: newFormCols.map(c => {
        const rules = c.is_required ? `:rules="[{ required: true, message: '请输入${c.comment}' }]"` : ''
        if (c.form_type === 'textarea') return `<a-form-item label="${c.comment}" name="${c.column_name}" ${rules}>\n  <a-textarea v-model:value="formData.${c.column_name}" />\n</a-form-item>`
        if (c.form_type === 'number') return `<a-form-item label="${c.comment}" name="${c.column_name}" ${rules}>\n  <a-input-number v-model:value="formData.${c.column_name}" style="width: 100%" />\n</a-form-item>`
        if (c.form_type === 'select') return `<a-form-item label="${c.comment}" name="${c.column_name}" ${rules}>\n  <a-select v-model:value="formData.${c.column_name}" />\n</a-form-item>`
        return `<a-form-item label="${c.comment}" name="${c.column_name}" ${rules}>\n  <a-input v-model:value="formData.${c.column_name}" />\n</a-form-item>`
      }).join('\n\n')
    })

    categories.push({
      key: 'form',
      label: `📝 表单变化 (${newFormCols.length})`,
      summary: [`开启表单显示：${newFormCols.map(c => c.comment || c.column_name).join('、')}`],
      guides
    })
  }

  // === 排序变化 Tab ===
  if (sortChanged) {
    const guides: ChangeGuide[] = []
    const listCols = currentCols.filter(c => c.column_name && c.is_list_visible)
    const formCols = currentCols.filter(c => c.column_name && c.is_form_visible)

    if (listCols.length > 0) {
      guides.push({
        title: `表格列新顺序`,
        description: `按新顺序调整 columns 数组`,
        file: `src/views/xxx/${moduleName}/index.vue`,
        code: listCols.map(c => `{ title: '${c.comment}', dataIndex: '${c.column_name}', key: '${c.column_name}' }`).join(',\n')
      })
    }

    if (formCols.length > 0) {
      guides.push({
        title: `表单项新顺序`,
        description: `按以下顺序调整表单项`,
        file: `src/views/xxx/${moduleName}/Form.vue`,
        code: `// 新的字段顺序:\n${formCols.map(c => c.column_name).join('\n')}`
      })
    }

    categories.push({
      key: 'sort',
      label: `🔀 排序变化`,
      summary: ['字段顺序已调整'],
      guides
    })
  }

  changeCategories.value = categories
  changeGuideTab.value = 'summary'
  changeGuideVisible.value = true
}

// 获取已生成模块
const fetchModules = async () => {
  try {
    const res = await getGeneratedModules()
    generatedModules.value = res.data
  } catch {
    // ignore
  }
}

// 删除已生成模块
const handleDeleteModule = async (name: string) => {
  await deleteModule(name)
  message.success('已删除生成代码')
  fetchModules()
}

// 获取已保存配置
const fetchSavedConfigs = async () => {
  configsLoading.value = true
  try {
    const res = await getSavedConfigs()
    savedConfigs.value = res.data || []
  } finally {
    configsLoading.value = false
  }
}

// 删除配置
const handleDeleteConfig = async (id: number) => {
  await deleteSavedConfig(id)
  message.success('配置已删除')
  fetchSavedConfigs()
}

// 获取菜单树
const fetchMenuTree = async () => {
  const res = await getMenuTree()
  menuTree.value = res.data
}

// 获取数据库表列表
const fetchDbTables = async () => {
  try {
    const res = await getTables()
    dbTables.value = res.data || []
  } catch {
    // ignore
  }
}

// 获取角色列表
const fetchRoleList = async () => {
  try {
    const res = await getRoleList()
    roleList.value = res.data || []
  } catch {
    // ignore
  }
}

// 获取字典类型列表
const fetchDictTypes = async () => {
  try {
    const res = await getAllDictTypes()
    dictTypeList.value = res.data || []
  } catch {
    // ignore
  }
}

onMounted(() => {
  fetchMenuTree()
  fetchModules()
  fetchSavedConfigs()
  fetchDbTables()
  fetchRoleList()
  fetchDictTypes()
})
</script>

<style scoped>
.generator-page {
  padding: 16px;
}
.code-preview {
  position: relative;
}
.code-toolbar {
  margin-bottom: 8px;
  padding: 8px 12px;
  background: #fafafa;
  border-radius: 4px 4px 0 0;
  border: 1px solid #e8e8e8;
  border-bottom: none;
}
.code-block {
  max-height: 550px;
  overflow: auto;
  margin: 0;
  padding: 16px;
  background: #1e1e1e;
  color: #d4d4d4;
  border-radius: 0 0 4px 4px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  white-space: pre;
  tab-size: 2;
}
.code-block code {
  font-family: inherit;
}
/* Go 语言高亮 */
.language-go {
  color: #9cdcfe;
}
/* TypeScript 高亮 */
.language-typescript {
  color: #4ec9b0;
}
/* Vue 高亮 */
.language-vue {
  color: #ce9178;
}
/* SQL 高亮 */
.language-sql {
  color: #dcdcaa;
}
.drawer-footer {
  position: absolute;
  right: 0;
  bottom: 0;
  width: 100%;
  border-top: 1px solid #e9e9e9;
  padding: 16px 24px;
  background: #fff;
  text-align: right;
  z-index: 1;
}
/* 拖拽排序样式 */
.draggable-table :deep(.drag-over) {
  background: #e6f7ff !important;
}
.draggable-table :deep(.dragging) {
  opacity: 0.5;
}
.draggable-table :deep(.drag-handle:hover) {
  color: #1890ff !important;
}
</style>
