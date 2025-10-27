package com.example.grpc.client.service;

import com.example.bookstore.*;
import com.example.grpc.client.model.AuthorDTO;
import com.example.grpc.client.model.BookDTO;
import com.example.grpc.client.model.ReviewDTO;
import com.example.grpc.client.util.TimeUtils;
import com.google.protobuf.Empty;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.client.inject.GrpcClient;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

@Service
@Slf4j
public class BookService {

    private final BookServiceGrpc.BookServiceBlockingStub stub;

    public BookService(@GrpcClient("book-server") BookServiceGrpc.BookServiceBlockingStub stub) {
        this.stub = stub;
    }

    public List<BookDTO.Response> findAll() {
        log.info("Fetching all books from gRPC server...");
        try {
            Iterator<Book> bookIterator = stub.findAll(Empty.newBuilder().build());
            List<BookDTO.Response> books = new ArrayList<>();
            while (bookIterator.hasNext()) {
                books.add(BookDTO.Response.from(bookIterator.next()));
            }
            log.info("Fetched {} books successfully", books.size());
            return books;
        } catch (Exception e) {
            log.error("Failed to fetch books: {}", e.getMessage(), e);
            throw e;
        }
    }

    public BookDTO.ResponseDetail get(Long id) {
        log.info("Fetching book with id={} from gRPC server...", id);
        try {
            BookDetail protoBook = stub.get(GetBookRequest.newBuilder().setId(id).build());
            log.info("Fetched book: id={}, title='{}'", protoBook.getId(), protoBook.getTitle());
            return BookDTO.ResponseDetail.from(protoBook,
                    protoBook.getAuthorsList().stream().map(AuthorDTO.Response::from).toList(),
                    protoBook.getReviewsList().stream().map(ReviewDTO.Response::from).toList());
        } catch (Exception e) {
            log.error("Failed to fetch book with id={}: {}", id, e.getMessage(), e);
            throw e;
        }
    }

    public BookDTO.Response create(BookDTO.Request request) {
        log.info("Creating new book with title='{}'...", request.getTitle());
        try {
            CreateBookRequest protoRequest = CreateBookRequest.newBuilder()
                    .setTitle(request.getTitle())
                    .setPublisher(request.getPublisher())
                    .setPublishedDate(TimeUtils.toTimestamp(request.getPublishedDate()))
                    .addAllAuthorIds(request.getAuthorIds())
                    .build();
            Book protoBook = stub.create(protoRequest);
            log.info("Book created: id={}, title='{}'", protoBook.getId(), protoBook.getTitle());
            return BookDTO.Response.from(protoBook);
        } catch (Exception e) {
            log.error("Failed to create book with title='{}': {}", request.getTitle(), e.getMessage(), e);
            throw e;
        }
    }

    public BookDTO.Response update(Long id, BookDTO.Request request) {
        log.info("Updating book id={} with new title='{}'...", id, request.getTitle());
        try {
            UpdateBookRequest protoRequest = UpdateBookRequest.newBuilder()
                    .setId(id)
                    .setTitle(request.getTitle())
                    .setPublisher(request.getPublisher())
                    .setPublishedDate(TimeUtils.toTimestamp(request.getPublishedDate()))
                    .addAllAuthorIds(request.getAuthorIds())
                    .build();
            Book protoBook = stub.update(protoRequest);
            log.info("Book updated: id={}, title='{}'", protoBook.getId(), protoBook.getTitle());
            return BookDTO.Response.from(protoBook);
        } catch (Exception e) {
            log.error("Failed to update book id={}: {}", id, e.getMessage(), e);
            throw e;
        }
    }

    public void delete(Long id) {
        log.info("Deleting book with id={}...", id);
        try {
            stub.delete(GetBookRequest.newBuilder().setId(id).build());
            log.info("Book with id={} deleted successfully", id);
        } catch (Exception e) {
            log.error("Failed to delete book id={}: {}", id, e.getMessage(), e);
            throw e;
        }
    }
}
