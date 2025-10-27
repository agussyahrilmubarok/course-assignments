package com.example.grpc.server.sevice;

import com.example.bookstore.*;
import com.example.grpc.server.domain.AuthorEntity;
import com.example.grpc.server.domain.BookEntity;
import com.example.grpc.server.events.BeforeDeleteAuthor;
import com.example.grpc.server.events.BeforeDeleteBook;
import com.example.grpc.server.repos.AuthorRepository;
import com.example.grpc.server.repos.BookRepository;
import com.example.grpc.server.repos.ReviewRepository;
import com.google.protobuf.Empty;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.server.service.GrpcService;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.context.event.EventListener;
import org.springframework.transaction.annotation.Transactional;

import java.time.*;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;

@GrpcService
@Slf4j
@Transactional
@RequiredArgsConstructor
public class BookServiceGrpcImpl extends BookServiceGrpc.BookServiceImplBase {

    private final BookRepository bookRepository;
    private final AuthorRepository authorRepository;
    private final ReviewRepository reviewRepository;
    private final ApplicationEventPublisher publisher;

    @Override
    public void findAll(Empty request, StreamObserver<Book> responseObserver) {
        Instant start = Instant.now();
        log.debug("BookService.findAll called");

        try {
            List<BookEntity> books = bookRepository.findAll();
            for (BookEntity book : books) {
                responseObserver.onNext(mapToBookProto(book));
            }
            responseObserver.onCompleted();
            log.info("findAll completed, total={}, duration={}ms",
                    books.size(), Duration.between(start, Instant.now()).toMillis());
        } catch (Exception e) {
            log.error("findAll failed", e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void get(GetBookRequest request, StreamObserver<BookDetail> responseObserver) {
        Instant start = Instant.now();
        log.debug("BookService.get called with id={}", request.getId());

        try {
            Optional<BookEntity> bookOpt = bookRepository.findById(request.getId());
            if (bookOpt.isPresent()) {
                responseObserver.onNext(mapToBookDetailProto(bookOpt.get()));
                responseObserver.onCompleted();
                log.info("get succeeded for id={}, duration={}ms",
                        request.getId(), Duration.between(start, Instant.now()).toMillis());
            } else {
                log.warn("get failed, book not found id={}", request.getId());
                responseObserver.onError(Status.NOT_FOUND
                        .withDescription("Book not found with id " + request.getId())
                        .asRuntimeException());
            }
        } catch (Exception e) {
            log.error("get failed for id={}", request.getId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    @Transactional
    public void create(CreateBookRequest request, StreamObserver<Book> responseObserver) {
        Instant start = Instant.now();
        log.debug("BookService.create called with title={}", request.getTitle());

        try {
            BookEntity book = new BookEntity();
            if (!request.getTitle().isEmpty()) {
                book.setTitle(request.getTitle());
            } else {
                responseObserver.onError(Status.INVALID_ARGUMENT.withDescription("Title is required").asRuntimeException());
                return;
            }
            book.setPublisher(request.getPublisher());
            book.setPublishedDate(fromTimestamp(request.getPublishedDate()));

            // Set authors
            Set<AuthorEntity> authors = request.getAuthorIdsList().stream()
                    .map(id -> authorRepository.findById(id).orElseThrow(() ->
                            Status.NOT_FOUND.withDescription("Author not found with id " + id).asRuntimeException()))
                    .collect(Collectors.toSet());
            book.setAuthors(authors);

            BookEntity saved = bookRepository.save(book);
            responseObserver.onNext(mapToBookProto(saved));
            responseObserver.onCompleted();
            log.info("create succeeded for id={}, duration={}ms",
                    saved.getId(), Duration.between(start, Instant.now()).toMillis());
        } catch (Exception e) {
            log.error("create failed", e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    @Transactional
    public void update(UpdateBookRequest request, StreamObserver<Book> responseObserver) {
        Instant start = Instant.now();
        log.debug("BookService.update called with id={}", request.getId());

        try {
            Optional<BookEntity> bookOpt = bookRepository.findById(request.getId());
            if (bookOpt.isEmpty()) {
                log.warn("update failed, book not found id={}", request.getId());
                responseObserver.onError(Status.NOT_FOUND
                        .withDescription("Book not found with id " + request.getId())
                        .asRuntimeException());
                return;
            }

            BookEntity book = bookOpt.get();

            if (!request.getTitle().isEmpty()) {
                book.setTitle(request.getTitle());
            }
            if (!request.getPublisher().isEmpty()) {
                book.setPublisher(request.getPublisher());
            }
            if (request.hasPublishedDate()) {
                book.setPublishedDate(fromTimestamp(request.getPublishedDate()));
            }

            // Update authors if provided
            if (!request.getAuthorIdsList().isEmpty()) {
                Set<AuthorEntity> authors = request.getAuthorIdsList().stream()
                        .map(id -> authorRepository.findById(id).orElseThrow(() ->
                                Status.NOT_FOUND.withDescription("Author not found with id " + id).asRuntimeException()))
                        .collect(Collectors.toSet());
                book.setAuthors(authors);
            }

            BookEntity updated = bookRepository.save(book);
            responseObserver.onNext(mapToBookProto(updated));
            responseObserver.onCompleted();
            log.info("update succeeded for id={}, duration={}ms",
                    updated.getId(), Duration.between(start, Instant.now()).toMillis());
        } catch (Exception e) {
            log.error("update failed for id={}", request.getId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    @Transactional
    public void delete(GetBookRequest request, StreamObserver<Empty> responseObserver) {
        Instant start = Instant.now();
        log.debug("BookService.delete called with id={}", request.getId());

        try {
            Optional<BookEntity> bookOpt = bookRepository.findById(request.getId());
            if (bookOpt.isEmpty()) {
                log.warn("delete failed, book not found id={}", request.getId());
                responseObserver.onError(Status.NOT_FOUND
                        .withDescription("Book not found with id " + request.getId())
                        .asRuntimeException());
                return;
            }

            BookEntity book = bookOpt.get();

            // Publish event before delete
            publisher.publishEvent(new BeforeDeleteBook(book.getId()));

            // Delete related reviews (cascade)
            reviewRepository.deleteAll(book.getReviews());

            // Delete book
            bookRepository.delete(book);

            responseObserver.onNext(Empty.getDefaultInstance());
            responseObserver.onCompleted();
            log.info("delete succeeded for id={}, duration={}ms",
                    request.getId(), Duration.between(start, Instant.now()).toMillis());
        } catch (Exception e) {
            log.error("delete failed for id={}", request.getId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    private LocalDate fromTimestamp(com.google.protobuf.Timestamp ts) {
        return (ts != null && ts.getSeconds() != 0)
                ? Instant.ofEpochSecond(ts.getSeconds(), ts.getNanos())
                .atZone(ZoneId.of("Asia/Jakarta"))
                .toLocalDate()
                : null;
    }


    private Book mapToBookProto(BookEntity entity) {
        Book.Builder builder = Book.newBuilder()
                .setId(entity.getId())
                .setTitle(entity.getTitle())
                .setPublisher(entity.getPublisher() != null ? entity.getPublisher() : "")
                .setPublishedDate(com.google.protobuf.Timestamp.newBuilder()
                        .setSeconds(entity.getPublishedDate() != null ?
                                entity.getPublishedDate().atStartOfDay(ZoneOffset.UTC).toEpochSecond() : 0)
                        .build())
                .setDateCreated(com.google.protobuf.Timestamp.newBuilder()
                        .setSeconds(entity.getDateCreated().toEpochSecond())
                        .setNanos(entity.getDateCreated().getNano())
                        .build())
                .setLastUpdated(com.google.protobuf.Timestamp.newBuilder()
                        .setSeconds(entity.getLastUpdated().toEpochSecond())
                        .setNanos(entity.getLastUpdated().getNano())
                        .build());

        entity.getAuthors().forEach(author -> builder.addAuthorIds(author.getId()));

        return builder.build();
    }

    private BookDetail mapToBookDetailProto(BookEntity entity) {
        BookDetail.Builder builder = BookDetail.newBuilder()
                .setId(entity.getId())
                .setTitle(entity.getTitle())
                .setPublisher(entity.getPublisher() != null ? entity.getPublisher() : "")
                .setPublishedDate(com.google.protobuf.Timestamp.newBuilder()
                        .setSeconds(entity.getPublishedDate() != null ?
                                entity.getPublishedDate().atStartOfDay(ZoneOffset.UTC).toEpochSecond() : 0)
                        .build())
                .setDateCreated(com.google.protobuf.Timestamp.newBuilder()
                        .setSeconds(entity.getDateCreated().toEpochSecond())
                        .setNanos(entity.getDateCreated().getNano())
                        .build())
                .setLastUpdated(com.google.protobuf.Timestamp.newBuilder()
                        .setSeconds(entity.getLastUpdated().toEpochSecond())
                        .setNanos(entity.getLastUpdated().getNano())
                        .build());

        // Map authors
        entity.getAuthors().forEach(author -> builder.addAuthors(Author.newBuilder()
                .setId(author.getId())
                .setName(author.getName())
                .setDateCreated(com.google.protobuf.Timestamp.newBuilder()
                        .setSeconds(author.getDateCreated().toEpochSecond())
                        .setNanos(author.getDateCreated().getNano())
                        .build())
                .setLastUpdated(com.google.protobuf.Timestamp.newBuilder()
                        .setSeconds(author.getLastUpdated().toEpochSecond())
                        .setNanos(author.getLastUpdated().getNano())
                        .build())
                .build()));

        // Map reviews
        entity.getReviews().forEach(review -> builder.addReviews(Review.newBuilder()
                .setId(review.getId())
                .setContent(review.getContent())
                .setRating(review.getRating().floatValue())
                .setBookId(entity.getId())
                .setDateCreated(com.google.protobuf.Timestamp.newBuilder()
                        .setSeconds(review.getDateCreated().toEpochSecond())
                        .setNanos(review.getDateCreated().getNano())
                        .build())
                .setLastUpdated(com.google.protobuf.Timestamp.newBuilder()
                        .setSeconds(review.getLastUpdated().toEpochSecond())
                        .setNanos(review.getLastUpdated().getNano())
                        .build())
                .build()));

        return builder.build();
    }

    @EventListener(BeforeDeleteAuthor.class)
    public void on(final BeforeDeleteAuthor event) {
        // remove many-to-many relations at owning side
        bookRepository.findAllByAuthors_Id(event.getId()).forEach(book ->
                book.getAuthors().removeIf(author -> author.getId().equals(event.getId())));
    }
}
