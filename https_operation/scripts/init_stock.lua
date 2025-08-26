-- KEYS[1]: 商品库存的 key (e.g., "product:stock:1")
-- KEYS[2]: 已秒杀用户集合的 key (e.g., "product:users:1")
-- ARGV[1]: 从 MySQL 读取的总库存数

if redis.call('EXISTS', KEYS[1]) == 1 then
    return 'already_initialized'
end
redis.call('SET', KEYS[1], ARGV[1])
redis.call('DEL', KEYS[2])
return 'ok'