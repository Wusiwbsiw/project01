-- KEYS[1]: 商品库存的 key (e.g., "product:stock:1")
-- KEYS[2]: 已秒杀用户集合的 key (e.g., "product:users:1")
-- ARGV[1]: 当前尝试秒杀的用户 ID

if redis.call('SISMEMBER', KEYS[2], ARGV[1]) == 1 then
  return 2 -- 2 代表“重复抢购”
end

local stock = tonumber(redis.call('GET', KEYS[1]))
if stock == nil or stock <= 0 then
  return 1 -- 1 代表“已售罄”
end

redis.call('DECR', KEYS[1])
redis.call('SADD', KEYS[2], ARGV[1])
return 0 -- 0 代表“秒杀成功”