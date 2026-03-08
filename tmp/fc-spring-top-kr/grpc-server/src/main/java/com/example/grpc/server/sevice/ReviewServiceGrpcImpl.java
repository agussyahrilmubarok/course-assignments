package com.example.grpc.server.sevice;

import com.example.bookstore.*;
import com.example.grpc.server.domain.ReviewEntity;
import com.example.grpc.server.events.BeforeDeleteBook;
import com.example.grpc.server.repos.BookRepository;
import com.example.grpc.server.repos.ReviewRepository;
import com.example.grpc.server.util.ReferencedException;
import com.google.protobuf.Empty;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.server.service.GrpcService;
import org.springframework.context.event.EventListener;
import org.springframework.transaction.annotation.Transactional;

import java.time.Duration;
import java.time.Instant;
import java.util.List;
import java.util.Optional;

@GrpcService
@Slf4j
@Transactional
@RequiredArgsConstructor
public class ReviewServiceGrpcImpl extends ReviewServiceGrpc.ReviewServiceImplBase {

    private final ReviewRepository reviewRepository;
    private final BookRepository bookRepository;

    @Override
    public void findAll(Empty request, StreamObserver<Review> responseObserver) {
        Instant start = Instant.now();
        log.debug("ReviewService.findAll called");

        try {
            List<ReviewEntity> reviews = reviewRepository.findAll();
            for (ReviewEntity entity : reviews) {
                responseObserver.onNext(mapToReviewProto(entity));
                log.debug("Review sent: id={}, bookId={}", entity.getId(), entity.getBook().getId());
            }
            responseObserver.onCompleted();
            log.info("findAll completed, total={}, duration={}ms",
                    reviews.size(), Duration.between(start, Instant.now()).toMillis());
        } catch (Exception e) {
            log.error("findAll failed", e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void get(GetReviewRequest request, StreamObserver<Review> responseObserver) {
        Instant start = Instant.now();
        log.debug("ReviewService.get called with id={}", request.getId());

        try {
            Optional<ReviewEntity> reviewOpt = reviewRepository.findById(request.getId());
            if (reviewOpt.isPresent()) {
                ReviewEntity entity = reviewOpt.get();
                responseObserver.onNext(mapToReviewProto(entity));
                responseObserver.onCompleted();
                log.info("get succeeded for id={}, duration={}ms",
                        request.getId(), Duration.between(start, Instant.now()).toMillis());
            } else {
                log.warn("get failed, review not found id={}", request.getId());
                responseObserver.onError(Status.NOT_FOUND
                        .withDescription("Review not found with id " + request.getId())
                        .asRuntimeException());
            }
        } catch (Exception e) {
            log.error("get failed for id={}", request.getId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    @Transactional
    public void create(CreateReviewRequest request, StreamObserver<Review> responseObserver) {
        Instant start = Instant.now();
        log.debug("ReviewService.create called with bookId={}, rating={}", request.getBookId(), request.getRating());

        try {
            var bookOpt = bookRepository.findById(request.getBookId());
            if (bookOpt.isEmpty()) {
                log.warn("create failed, book not found id={}", request.getBookId());
                responseObserver.onError(Status.NOT_FOUND
                        .withDescription("Book not found with id " + request.getBookId())
                        .asRuntimeException());
                return;
            }

            ReviewEntity entity = new ReviewEntity();
            entity.setContent(request.getContent() != null ? request.getContent() : "");
            entity.setRating(request.getRating() != 0 ? (double) request.getRating() : 0.0);
            entity.setBook(bookOpt.get());

            ReviewEntity saved = reviewRepository.save(entity);
            responseObserver.onNext(mapToReviewProto(saved));
            responseObserver.onCompleted();
            log.info("create succeeded for id={}, duration={}ms",
                    saved.getId(), Duration.between(start, Instant.now()).toMillis());
        } catch (Exception e) {
            log.error("create failed for bookId={}", request.getBookId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    @Transactional
    public void update(UpdateReviewRequest request, StreamObserver<Review> responseObserver) {
        Instant start = Instant.now();
        log.debug("ReviewService.update called with id={}", request.getId());

        try {
            Optional<ReviewEntity> reviewOpt = reviewRepository.findById(request.getId());
            if (reviewOpt.isEmpty()) {
                log.warn("update failed, review not found id={}", request.getId());
                responseObserver.onError(Status.NOT_FOUND
                        .withDescription("Review not found with id " + request.getId())
                        .asRuntimeException());
                return;
            }

            ReviewEntity entity = reviewOpt.get();

            if (request.getBookId() != 0) {
                var bookOpt = bookRepository.findById(request.getBookId());
                if (bookOpt.isEmpty()) {
                    log.warn("update failed, book not found id={}", request.getBookId());
                    responseObserver.onError(Status.NOT_FOUND
                            .withDescription("Book not found with id " + request.getBookId())
                            .asRuntimeException());
                    return;
                }
                entity.setBook(bookOpt.get());
            }

            entity.setContent(request.getContent() != null ? request.getContent() : entity.getContent());
            entity.setRating(request.getRating() != 0 ? (double) request.getRating() : entity.getRating());

            ReviewEntity updated = reviewRepository.save(entity);
            responseObserver.onNext(mapToReviewProto(updated));
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
    public void delete(GetReviewRequest request, StreamObserver<Empty> responseObserver) {
        Instant start = Instant.now();
        log.debug("ReviewService.delete called with id={}", request.getId());

        try {
            Optional<ReviewEntity> reviewOpt = reviewRepository.findById(request.getId());
            if (reviewOpt.isEmpty()) {
                log.warn("delete failed, review not found id={}", request.getId());
                responseObserver.onError(Status.NOT_FOUND
                        .withDescription("Review not found with id " + request.getId())
                        .asRuntimeException());
                return;
            }

            reviewRepository.delete(reviewOpt.get());
            responseObserver.onNext(Empty.getDefaultInstance());
            responseObserver.onCompleted();
            log.info("delete succeeded for id={}, duration={}ms",
                    request.getId(), Duration.between(start, Instant.now()).toMillis());
        } catch (Exception e) {
            log.error("delete failed for id={}", request.getId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    private Review mapToReviewProto(ReviewEntity entity) {
        return Review.newBuilder()
                .setId(entity.getId())
                .setContent(entity.getContent())
                .setRating(entity.getRating() != null ? entity.getRating().floatValue() : 0f)
                .setBookId(entity.getBook().getId())
                .setDateCreated(com.google.protobuf.Timestamp.newBuilder()
                        .setSeconds(entity.getDateCreated().toEpochSecond())
                        .setNanos(entity.getDateCreated().getNano())
                        .build())
                .setLastUpdated(com.google.protobuf.Timestamp.newBuilder()
                        .setSeconds(entity.getLastUpdated().toEpochSecond())
                        .setNanos(entity.getLastUpdated().getNano())
                        .build())
                .build();
    }

    @EventListener(BeforeDeleteBook.class)
    public void on(final BeforeDeleteBook event) {
        ReviewEntity bookReview = reviewRepository.findFirstByBook_Id(event.getId());
        if (bookReview != null) {
            ReferencedException referencedException = new ReferencedException();
            referencedException.setKey("book.review.book.referenced");
            referencedException.addParam(bookReview.getId());
            throw referencedException;
        }
    }
}
