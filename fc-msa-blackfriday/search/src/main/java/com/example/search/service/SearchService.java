package com.example.search.service;

import java.util.List;

public interface SearchService {

    List<String> findProductIdsByTag(String tag);

    void addTagsCache(String productId, List<String> tags);

    void removeTagsCache(String productId, List<String> tags);
}
