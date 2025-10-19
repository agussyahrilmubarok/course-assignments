package com.example.rest.model;

import com.example.rest.domain.Author;
import com.fasterxml.jackson.annotation.JsonInclude;
import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.*;

import java.time.OffsetDateTime;


@Getter
@Setter
public class AuthorDTO {

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    @Schema(name = "AuthorRequest")
    public static class Request {
        @NotNull
        @Size(max = 255)
        private String name;
    }

    @Getter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @JsonInclude(JsonInclude.Include.ALWAYS)
    @Schema(name = "AuthorResponse")
    public static class Response {

        private Long id;
        private String name;
        private OffsetDateTime dateCreated;
        private OffsetDateTime lastUpdated;

        public static Response from(Author author) {
            return Response.builder()
                    .id(author.getId())
                    .name(author.getName())
                    .dateCreated(author.getDateCreated())
                    .lastUpdated(author.getLastUpdated())
                    .build();
        }
    }
}
