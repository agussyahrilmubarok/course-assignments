package com.example.grpc.server.sevice;

import com.example.bookstore.*;
import com.example.grpc.server.domain.AuthorEntity;
import com.example.grpc.server.events.BeforeDeleteAuthor;
import com.example.grpc.server.repos.AuthorRepository;
import com.google.protobuf.Empty;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.server.service.GrpcService;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.data.domain.Sort;
import org.springframework.transaction.annotation.Transactional;

import java.time.Duration;
import java.time.Instant;
import java.util.List;
import java.util.Optional;

@GrpcService
@Slf4j
@Transactional
@RequiredArgsConstructor
public class AuthorServiceGrpcImpl extends AuthorServiceGrpc.AuthorServiceImplBase {

    private final AuthorRepository authorRepository;
    private final ApplicationEventPublisher publisher;

    @Override
    public void findAll(Empty request, StreamObserver<Author> responseObserver) {
        Instant start = Instant.now();
        log.debug("AuthorService.findAll called");

        try {
            List<AuthorEntity> authors = authorRepository.findAll(Sort.by("id"));
            authors.forEach(entity -> {
                responseObserver.onNext(mapToAuthorProto(entity));
                log.debug("Author sent: id={}, name={}", entity.getId(), entity.getName());
            });
            responseObserver.onCompleted();
            log.info("findAll completed, total={}, duration={}ms",
                    authors.size(), Duration.between(start, Instant.now()).toMillis());
        } catch (Exception e) {
            log.error("findAll failed", e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void get(GetAuthorRequest request, StreamObserver<Author> responseObserver) {
        Instant start = Instant.now();
        log.debug("AuthorService.get called with id={}", request.getId());

        try {
            Optional<AuthorEntity> authorOpt = authorRepository.findById(request.getId());
            if (authorOpt.isPresent()) {
                AuthorEntity entity = authorOpt.get();
                responseObserver.onNext(mapToAuthorProto(entity));
                responseObserver.onCompleted();
                log.info("get succeeded for id={}, duration={}ms", request.getId(), Duration.between(start, Instant.now()).toMillis());
            } else {
                log.warn("get failed, author not found id={}", request.getId());
                responseObserver.onError(Status.NOT_FOUND
                        .withDescription("Author not found with id " + request.getId())
                        .asRuntimeException());
            }
        } catch (Exception e) {
            log.error("get failed for id={}", request.getId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void create(CreateAuthorRequest request, StreamObserver<Author> responseObserver) {
        Instant start = Instant.now();
        log.debug("AuthorService.create called with name={}", request.getName());

        try {
            AuthorEntity entity = new AuthorEntity();
            entity.setName(request.getName());
            AuthorEntity saved = authorRepository.save(entity);

            responseObserver.onNext(mapToAuthorProto(saved));
            responseObserver.onCompleted();
            log.info("create succeeded for id={}, duration={}ms", saved.getId(), Duration.between(start, Instant.now()).toMillis());
        } catch (Exception e) {
            log.error("create failed for name={}", request.getName(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void update(UpdateAuthorRequest request, StreamObserver<Author> responseObserver) {
        Instant start = Instant.now();
        log.debug("AuthorService.update called with id={}", request.getId());

        try {
            Optional<AuthorEntity> authorOpt = authorRepository.findById(request.getId());
            if (authorOpt.isPresent()) {
                AuthorEntity entity = authorOpt.get();
                entity.setName(request.getName());
                AuthorEntity updated = authorRepository.save(entity);

                responseObserver.onNext(mapToAuthorProto(updated));
                responseObserver.onCompleted();
                log.info("update succeeded for id={}, duration={}ms", updated.getId(), Duration.between(start, Instant.now()).toMillis());
            } else {
                log.warn("update failed, author not found id={}", request.getId());
                responseObserver.onError(Status.NOT_FOUND
                        .withDescription("Author not found with id " + request.getId())
                        .asRuntimeException());
            }
        } catch (Exception e) {
            log.error("update failed for id={}", request.getId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void delete(GetAuthorRequest request, StreamObserver<Empty> responseObserver) {
        Instant start = Instant.now();
        log.debug("AuthorService.delete called with id={}", request.getId());

        try {
            Optional<AuthorEntity> authorOpt = authorRepository.findById(request.getId());
            if (authorOpt.isPresent()) {
                publisher.publishEvent(new BeforeDeleteAuthor(request.getId()));
                authorRepository.delete(authorOpt.get());

                responseObserver.onNext(Empty.getDefaultInstance());
                responseObserver.onCompleted();
                log.info("delete succeeded for id={}, duration={}ms", request.getId(), Duration.between(start, Instant.now()).toMillis());
            } else {
                log.warn("delete failed, author not found id={}", request.getId());
                responseObserver.onError(Status.NOT_FOUND
                        .withDescription("Author not found with id " + request.getId())
                        .asRuntimeException());
            }
        } catch (Exception e) {
            log.error("delete failed for id={}", request.getId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    private Author mapToAuthorProto(AuthorEntity entity) {
        return Author.newBuilder()
                .setId(entity.getId())
                .setName(entity.getName())
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
}
