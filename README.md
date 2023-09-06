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
![1](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/bb6a19f5-24bc-4cfc-851a-e5af7dfc72d7)
![2](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/80b8bbb1-512c-45da-bca6-29609dd5e0d0)
![3](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/9962cc74-622c-4b7a-b397-9c78e8f10c8a)
![4](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/371ffdd4-032e-4f0c-9c56-af7044aeca5b)
![5](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/c7e74981-a90b-4acf-9a4d-938a20d87dac)
![6](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/656f90b5-eb7d-49e4-8443-0fa54f90f5a9)
![7](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/536f6174-9455-4ed6-b83c-aaec660af46b)

#### 添加作者
![WechatIMG529](https://github.com/Mlegbder/midjourney-wechat/assets/28382910/84031f9e-5b25-4898-89e4-4f030ba1f443)
