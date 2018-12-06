--
-- Created by IntelliJ IDEA.
-- User: ace
-- Date: 2018/4/16
-- Time: 下午11:34
-- To change this template use File | Settings | File Templates.
--

-- cache stock
local stock = KEYS[1]
local stock_lock = KEYS[2]
local stock_val = tonumber(ARGV[1])
local stock_lock_val = tonumber(ARGV[2])
local is_exists = redis.call('EXISTS', stock)
if is_exists == 1 then
    return 0
else
    redis.call('set', stock, stock_val)
    redis.call('set', stock_lock, stock_lock_val)
    return 1
end


-- reduce stock  0 key不存在 -1 库存不足 返回list 成功
local stock = KEYS[1]
local stock_lock = KEYS[2]
local stock_change = tonumber(ARGV[1])
local stock_lock_change = tonumber(ARGV[2])
local is_exists = redis.call('EXISTS', stock)
if is_exists == 1 then
    local stockAftChange = redis.call("incrby", stock, stock_change);
    if(stockAftChange < 0) then
        redis.call("decrby", stock, stock_change);
        return -1
    else
        local stockLockAftChange = redis.call("incrby", stock_lock, stock_lock_change);
        return {stockAftChange,stockLockAftChange};
    end
else
    return 0;
end