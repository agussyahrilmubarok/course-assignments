package com.example.order.model;

import com.fasterxml.jackson.annotation.JsonInclude;
import lombok.*;

import java.util.List;

@Getter
@Setter
public class ProductDTO {

    @Getter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @JsonInclude(JsonInclude.Include.ALWAYS)
    public static class Response {
        private String id;
        private String sellerId;
        private String name;
        private String description;
        private Long price;
        private Long stockCount;
        private List<String> tags;
    }
}
