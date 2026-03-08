package com.example.graphql.service;

import com.example.graphql.domain.Author;
import com.example.graphql.events.BeforeDeleteAuthor;
import com.example.graphql.model.AuthorDTO;
import com.example.graphql.repos.AuthorRepository;
import com.example.graphql.util.NotFoundException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
@Slf4j
@Transactional(rollbackFor = Exception.class)
public class AuthorService {

    private final AuthorRepository authorRepository;
    private final ApplicationEventPublisher publisher;

    public AuthorService(final AuthorRepository authorRepository,
                         final ApplicationEventPublisher publisher) {
        this.authorRepository = authorRepository;
        this.publisher = publisher;
    }

    public List<AuthorDTO.Response> findAll() {
        log.info("Fetching all authors");
        final List<Author> authors = authorRepository.findAll(Sort.by("id"));
        log.debug("Found {} authors", authors.size());

        return authors.stream()
                .map(AuthorDTO.Response::from)
                .toList();
    }

    public AuthorDTO.Response get(final Long id) {
        log.info("Fetching author with id: {}", id);
        return authorRepository.findById(id)
                .map(author -> {
                    log.debug("Author found: {}", author.getName());
                    return AuthorDTO.Response.from(author);
                })
                .orElseThrow(() -> {
                    log.error("Author with id {} not found", id);
                    return new NotFoundException();
                });
    }

    public AuthorDTO.Response create(final AuthorDTO.AuthorRequest param) {
        log.info("Creating new author with name: {}", param.getName());
        final Author author = new Author();
        author.setName(param.getName());
        final Author saved = authorRepository.save(author);
        log.debug("Author created with id: {}", saved.getId());

        return AuthorDTO.Response.from(saved);
    }

    public AuthorDTO.Response update(final Long id, final AuthorDTO.AuthorRequest param) {
        log.info("Updating author with id: {}", id);
        final Author author = authorRepository.findById(id)
                .orElseThrow(() -> {
                    log.error("Author with id {} not found for update", id);
                    return new NotFoundException();
                });

        author.setName(param.getName());
        final Author updated = authorRepository.save(author);
        log.debug("Author with id {} updated to name: {}", id, updated.getName());

        return AuthorDTO.Response.from(updated);
    }

    public void delete(final Long id) {
        log.info("Deleting author with id: {}", id);
        final Author author = authorRepository.findById(id)
                .orElseThrow(() -> {
                    log.error("Author with id {} not found for deletion", id);
                    return new NotFoundException();
                });

        publisher.publishEvent(new BeforeDeleteAuthor(id));
        authorRepository.delete(author);
        log.debug("Author with id {} deleted successfully", id);
    }
}
