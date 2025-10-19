package com.example.search.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RSet;
import org.redisson.api.RedissonClient;
import org.springframework.stereotype.Service;

import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

@Service
@Slf4j
@RequiredArgsConstructor
public class SearchServiceImpl implements SearchService {

    private static final String PRODUCT_TAG_PREFIX = "product-tags:";

    private final RedissonClient redissonClient;

    @Override
    public List<String> findProductIdsByTag(String tag) {
        String redisKey = PRODUCT_TAG_PREFIX + tag;
        RSet<String> productSet = redissonClient.getSet(redisKey);
        if (productSet == null || productSet.isEmpty()) {
            log.info("No product IDs found for tag '{}'", tag);
            return Collections.emptyList();
        }

        List<String> productIds = productSet.stream().collect(Collectors.toList());
        log.info("Found {} product IDs for tag '{}': {}", productIds.size(), tag, productIds);
        return productIds;
    }

    @Override
    public void addTagsCache(String productId, List<String> tags) {
        tags.forEach(tag -> {
            String redisKey = PRODUCT_TAG_PREFIX + tag;
            RSet<String> productSet = redissonClient.getSet(redisKey);
            boolean added = productSet.add(productId);
            log.debug("🔹 Redis key '{}': productId '{}' added = {}", redisKey, productId, added);
        });
        log.info("Finished adding productId '{}' to tags", productId);
    }

    @Override
    public void removeTagsCache(String productId, List<String> tags) {
        tags.forEach(tag -> {
            String redisKey = PRODUCT_TAG_PREFIX + tag;
            RSet<String> productSet = redissonClient.getSet(redisKey);
            boolean removed = productSet.remove(productId);
            log.debug("Redis key '{}': productId '{}' removed = {}", redisKey, productId, removed);
        });
        log.info("Finished removing productId '{}' from tags", productId);
    }
}
