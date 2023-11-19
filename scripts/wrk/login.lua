wrk.method="POST"
wrk.headers["Content-Type"] = "application/json"
-- 这个要改为你的注册的数据
wrk.body='{"email":"837267496@qq.com", "password": "12210917cc!"}'

-- 运行：wrk -t4 -d30s -c2 -s ./scripts/wrk/login.lua http://localhost:8081/user/login

