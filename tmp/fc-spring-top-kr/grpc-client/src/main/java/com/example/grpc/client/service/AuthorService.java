package com.example.grpc.client.service;

import com.example.bookstore.*;
import com.example.grpc.client.model.AuthorDTO;
import com.google.protobuf.Empty;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.client.inject.GrpcClient;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

@Service
@Slf4j
public class AuthorService {

    private final AuthorServiceGrpc.AuthorServiceBlockingStub stub;

    public AuthorService(@GrpcClient("book-server") AuthorServiceGrpc.AuthorServiceBlockingStub stub) {
        this.stub = stub;
    }

    public List<AuthorDTO.Response> findAll() {
        log.info("Fetching all authors from gRPC server...");

        try {
            Iterator<Author> authorIterator = stub.findAll(Empty.newBuilder().build());
            List<AuthorDTO.Response> authors = new ArrayList<>();
            while (authorIterator.hasNext()) {
                authors.add(AuthorDTO.Response.from(authorIterator.next()));
            }
            log.info("Fetched {} authors successfully", authors.size());
            return authors;
        } catch (Exception e) {
            log.error("Failed to fetch authors: {}", e.getMessage(), e);
            throw e;
        }
    }

    public AuthorDTO.Response get(Long id) {
        log.info("Fetching author with id={} from gRPC server...", id);
        try {
            Author protoAuthor = stub.get(GetAuthorRequest.newBuilder().setId(id).build());
            log.info("Fetched author: id={}, name={}", protoAuthor.getId(), protoAuthor.getName());
            return AuthorDTO.Response.from(protoAuthor);
        } catch (Exception e) {
            log.error("Failed to fetch author with id={}: {}", id, e.getMessage(), e);
            throw e;
        }
    }

    public AuthorDTO.Response create(AuthorDTO.Request request) {
        log.info("Creating new author with name='{}'...", request.getName());
        try {
            var protoRequest = CreateAuthorRequest.newBuilder()
                    .setName(request.getName())
                    .build();
            Author protoAuthor = stub.create(protoRequest);
            log.info("Author created: id={}, name={}", protoAuthor.getId(), protoAuthor.getName());
            return AuthorDTO.Response.from(protoAuthor);
        } catch (Exception e) {
            log.error("Failed to create author with name='{}': {}", request.getName(), e.getMessage(), e);
            throw e;
        }
    }

    public AuthorDTO.Response update(Long id, AuthorDTO.Request request) {
        log.info("Updating author id={} with new name='{}'...", id, request.getName());
        try {
            var protoRequest = UpdateAuthorRequest.newBuilder()
                    .setId(id)
                    .setName(request.getName())
                    .build();
            Author protoAuthor = stub.update(protoRequest);
            log.info("Author updated: id={}, name={}", protoAuthor.getId(), protoAuthor.getName());
            return AuthorDTO.Response.from(protoAuthor);
        } catch (Exception e) {
            log.error("Failed to update author id={}: {}", id, e.getMessage(), e);
            throw e;
        }
    }

    public void delete(Long id) {
        log.info("Deleting author with id={}...", id);
        try {
            stub.delete(GetAuthorRequest.newBuilder().setId(id).build());
            log.info("Author with id={} deleted successfully", id);
        } catch (Exception e) {
            log.error("Failed to delete author with id={}: {}", id, e.getMessage(), e);
            throw e;
        }
    }
}
