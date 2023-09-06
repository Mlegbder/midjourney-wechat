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
- /imagine [图片链接(非必填)] 1gril[描述词]
- /up [任务ID] U1 , /up [任务ID] V1
- /help(查看帮助)
- /cz 2 100 , /cz [用户id] [余额]
- /me(查个人ID及余额)
- /describe [图片链接]
- /blend [图片链接] [图片链接]

#### 如果无需充值功能的,屏蔽mysql初始化代码即可
 
#### 功能截图
![1](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/6f3a41f8-9304-4f9d-93c9-d974dfd2f1c8)
![3](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/66a13b61-da9f-4559-bbc9-cfd092178c44)
![5](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/2f1f2827-30fa-4401-a096-45421966e66f)
![6](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/035b08c8-f831-49c5-b327-9c20809dcd03)
![7](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/7c469e76-13ef-42b3-a9d6-d16399b95ff3)


#### 添加作者
![WechatIMG529](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/d636bb09-65c1-450c-9873-7f239fd45f9d)
