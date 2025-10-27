package com.example.rest.service;

import com.example.rest.domain.Author;
import com.example.rest.events.BeforeDeleteAuthor;
import com.example.rest.model.AuthorDTO;
import com.example.rest.repos.AuthorRepository;
import com.example.rest.util.NotFoundException;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;


@Service
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
        final List<Author> authors = authorRepository.findAll(Sort.by("id"));

        return authors.stream()
                .map(AuthorDTO.Response::from)
                .toList();
    }

    public AuthorDTO.Response get(final Long id) {
        return authorRepository.findById(id)
                .map(AuthorDTO.Response::from)
                .orElseThrow(NotFoundException::new);
    }

    public AuthorDTO.Response create(final AuthorDTO.Request param) {
        final Author author = new Author();
        author.setName(param.getName());

        return AuthorDTO.Response.from(authorRepository.save(author));
    }

    public AuthorDTO.Response update(final Long id, final AuthorDTO.Request param) {
        final Author author = authorRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        author.setName(param.getName());

        return AuthorDTO.Response.from(authorRepository.save(author));
    }

    public void delete(final Long id) {
        final Author author = authorRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        publisher.publishEvent(new BeforeDeleteAuthor(id));
        authorRepository.delete(author);
    }
}
