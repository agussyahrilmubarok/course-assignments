package com.example.search.rest;

import com.example.search.service.SearchService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController("SearchResourceV1")
@RequestMapping(value = "/api/v1/search", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class SearchResource {

    private final SearchService searchService;

    @GetMapping("/tags/{tag}/productIds")
    public ResponseEntity<List<String>> getTagProductIds(@PathVariable String tag) {
        return ResponseEntity.ok(searchService.findProductIdsByTag(tag));
    }
}
