package com.example.graphql.model;

import com.example.graphql.domain.Review;
import com.fasterxml.jackson.annotation.JsonInclude;
import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.NotNull;
import lombok.*;

import java.time.OffsetDateTime;


@Getter
@Setter
public class ReviewDTO {

    @Data
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    @Schema(name = "ReviewRequest")
    public static class ReviewRequest {

        @NotNull
        private String content;

        private Double rating;

        @NotNull
        private Long bookId;
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @JsonInclude(JsonInclude.Include.ALWAYS)
    @Schema(name = "ReviewResponse")
    public static class Response {

        private Long id;
        private String content;
        private Double rating;
        private Long bookId;
        private OffsetDateTime dateCreated;
        private OffsetDateTime lastUpdated;

        public static ReviewDTO.Response from(Review review) {
            return Response.builder()
                    .id(review.getId())
                    .content(review.getContent())
                    .rating(review.getRating())
                    .bookId(review.getBook() == null ? null : review.getBook().getId())
                    .dateCreated(review.getDateCreated())
                    .lastUpdated(review.getLastUpdated())
                    .build();
        }
    }
}