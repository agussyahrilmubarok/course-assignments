package com.example.search.dto;

import lombok.*;

import java.util.List;

@Getter
@Setter
public class ProductDTO {

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    public static class ProductTagsMessage {

        public String productId;

        public List<String> tags;
    }
}
