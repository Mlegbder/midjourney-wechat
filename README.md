# midjourney-wechat

#### 微信画图机器人

#### 该项目基于
  - [midjourney-proxy](https://github.com/novicezk/midjourney-proxy) - 自行按照教程部署(调用此代理服务)
  - [midjourney-wechat](https://github.com/geek-bigniu/midjourney-wechat) - 基于此项目修改

#### 功能实现
- [x] 支持 Imagine 指令和相关动作
- [x] 支持 Blend(图片混合)
- [x] 支持 Describe(图生文) 指令
- [x] 支持 放大(U), 变换(V) 指令
- [x] 自动注册, 查询个人余额
- [x] 管理员充值命令

#### 指令示例
/imagine(画图)
  -- /imagine [图片链接(非必填)] 1gril[描述词]
/up(变换,放大)
 -- /up [任务ID] U1 , /up [任务ID] V1
/help(查看帮助)
/cz(充值)
 -- /cz 2 100 , /cz [用户id] [余额]
/me(查个人ID及余额)
/describe(反推图升文)
 -- /describe [图片链接]
/blend(混图)
 -- /blend [图片链接] [图片链接]
 
