
在此处记录立项过程以及具体如何使用

# project01
gogogo

测试主页


测试登陆
curl -k -X POST \
  -H "Content-Type: application/json" \
  -d '{"userid": 100000, "password": "123456"}' \
  https://localhost:8443/api/user/login

测试注册
curl -k -X POST \
      -H "Content-Type: application/json" \
      -d '{"username": "alice", "password": "alice_password_123"}' \
      https://localhost:8443/api/user/register

测试修改昵称
curl -k -X POST \
  -H "Content-Type: application/json" \
  -d '{"userid": 100000, "newname": "NewAwesomeName"}' \
  https://localhost:8443/api/user/reset-name

测试修改密码
curl -k -X POST \
  -H "Content-Type: application/json" \
  -d '{"userid": 100000, "newpassword": "MyNewSecurePassword!@#"}' \
  https://localhost:8443/api/user/reset-password

