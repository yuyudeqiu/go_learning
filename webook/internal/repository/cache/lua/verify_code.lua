local key = KEYS[1]
local cntKey = key .. ":cnt"

local expectedCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)

if cnt == nil or cnt <= 0 then
    --    验证次数耗尽了
    return -1
end

if code == expectedCode then -- 不相等，用户输错了
    redis.call("set", cntKey, 0)
    return 0
else

    redis.call("decr", cntKey)

    return -2
end