## 本项目简单介绍

1. 管理模式： git flow 的特性分支开发
- master 分支：用于发布版本
- develop 分支：用于日常开发
- feature 分支：用于开发新功能
- hotfix 分支：用于修复 bug

2. 命名规范：
- 分支命名：feature/xxx、hotfix/xxx
- 只有短期分支需要有层级
- feature 分支命名：feature/redPacket
- 如果有必要，feature分支可以延伸第三层级，但是最好不要有第四层级，比如：feature/redPacket/send

3.commit规范：
- feat(cloud_wallet): 新增红包功能 
- 上面的部分可以拆成3部分来看
  - feat 代表这次提交的类型，可以是 feat、fix、docs、style、refactor、test、chore
  - cloud_wallet 代表这次提交的服务，根据自身服务来命名，比如msg
  - 具体描述：需要清晰的具体描述这次提交的内容，比如新增红包功能