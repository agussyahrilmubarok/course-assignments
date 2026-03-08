package com.example.grpc.client.model;

import com.example.bookstore.Book;
import com.example.bookstore.BookDetail;
import com.example.grpc.client.util.TimeUtils;
import com.fasterxml.jackson.annotation.JsonInclude;
import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.*;

import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.List;


@Getter
@Setter
public class BookDTO {

    @Getter
    @Builder
    @NoArgsConstructor(access = AccessLevel.PRIVATE)
    @AllArgsConstructor(access = AccessLevel.PRIVATE)
    @Schema(name = "BookRequest")
    public static class Request {

        @NotNull
        @Size(max = 255)
        private String title;

        @Size(max = 255)
        private String publisher;

        @NotNull
        private LocalDate publishedDate;

        private List<Long> authorIds;
    }

    @Getter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @JsonInclude(JsonInclude.Include.ALWAYS)
    @Schema(name = "BookResponse")
    public static class Response {

        private Long id;
        private String title;
        private String publisher;
        private LocalDate publishedDate;
        private OffsetDateTime dateCreated;
        private OffsetDateTime lastUpdated;

        public static BookDTO.Response from(Book book) {
            return Response.builder()
                    .id(book.getId())
                    .title(book.getTitle())
                    .publisher(book.getPublisher())
                    .publishedDate(TimeUtils.tsToLocalDate(book.getPublishedDate()))
                    .dateCreated(TimeUtils.tsToOffsetDateTime(book.getDateCreated()))
                    .lastUpdated(TimeUtils.tsToOffsetDateTime(book.getLastUpdated()))
                    .build();
        }
    }

    @Getter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @JsonInclude(JsonInclude.Include.ALWAYS)
    @Schema(name = "BookResponse")
    public static class ResponseDetail {

        private Long id;
        private String title;
        private String publisher;
        private LocalDate publishedDate;
        private List<AuthorDTO.Response> authors;
        private List<ReviewDTO.Response> reviews;
        private OffsetDateTime dateCreated;
        private OffsetDateTime lastUpdated;

        public static BookDTO.ResponseDetail from(BookDetail book, List<AuthorDTO.Response> authors, List<ReviewDTO.Response> reviews) {
            return ResponseDetail.builder()
                    .id(book.getId())
                    .title(book.getTitle())
                    .publisher(book.getPublisher())
                    .publishedDate(TimeUtils.tsToLocalDate(book.getPublishedDate()))
                    .authors(authors)
                    .reviews(reviews)
                    .dateCreated(TimeUtils.tsToOffsetDateTime(book.getDateCreated()))
                    .lastUpdated(TimeUtils.tsToOffsetDateTime(book.getLastUpdated()))
                    .build();
        }
    }
}