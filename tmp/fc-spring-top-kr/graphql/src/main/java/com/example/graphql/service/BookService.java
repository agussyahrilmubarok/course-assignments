package com.example.graphql.service;

import com.example.graphql.domain.Author;
import com.example.graphql.domain.Book;
import com.example.graphql.events.BeforeDeleteAuthor;
import com.example.graphql.events.BeforeDeleteBook;
import com.example.graphql.model.AuthorDTO;
import com.example.graphql.model.BookDTO;
import com.example.graphql.model.ReviewDTO;
import com.example.graphql.repos.AuthorRepository;
import com.example.graphql.repos.BookRepository;
import com.example.graphql.repos.ReviewRepository;
import com.example.graphql.util.NotFoundException;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.context.event.EventListener;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.HashSet;
import java.util.List;


@Service
@Transactional(rollbackFor = Exception.class)
public class BookService {

    private final BookRepository bookRepository;
    private final AuthorRepository authorRepository;
    private final ReviewRepository reviewRepository;
    private final ApplicationEventPublisher publisher;

    public BookService(final BookRepository bookRepository, final AuthorRepository authorRepository, ReviewRepository reviewRepository,
                       final ApplicationEventPublisher publisher) {
        this.bookRepository = bookRepository;
        this.authorRepository = authorRepository;
        this.reviewRepository = reviewRepository;
        this.publisher = publisher;
    }

    public List<BookDTO.Response> findAll() {
        final List<Book> books = bookRepository.findAll(Sort.by("id"));

        return books.stream()
                .map(book -> BookDTO.Response.from(book,
                        book.getAuthors().stream().map(AuthorDTO.Response::from).toList(),
                        book.getReviews().stream().map(ReviewDTO.Response::from).toList()))
                .toList();
    }

    @Transactional
    public BookDTO.Response get(final Long id) {
        return bookRepository.findById(id)
                .map(book -> BookDTO.Response.from(book,
                        book.getAuthors().stream().map(AuthorDTO.Response::from).toList(),
                        book.getReviews().stream().map(ReviewDTO.Response::from).toList()))
                .orElseThrow(NotFoundException::new);
    }

    @Transactional
    public BookDTO.Response create(final BookDTO.BookRequest param) {
        Book book = new Book();
        book.setTitle(param.getTitle());
        book.setPublisher(param.getPublisher());
        book.setPublishedDate(param.getPublishedDate());
        final List<Author> authors = authorRepository.findAllById(param.getAuthorIds() == null ? List.of() : param.getAuthorIds());
        if (authors.size() != (param.getAuthorIds() == null ? 0 : param.getAuthorIds().size())) {
            throw new NotFoundException("one of authors not found");
        }
        book.setAuthors(new HashSet<>(authors));
        book = bookRepository.save(book);

        return BookDTO.Response.from(book,
                book.getAuthors().stream().map(AuthorDTO.Response::from).toList(),
                book.getReviews().stream().map(ReviewDTO.Response::from).toList());
    }

    @Transactional
    public BookDTO.Response update(final Long id, final BookDTO.BookRequest param) {
        Book book = bookRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        book.setTitle(param.getTitle());
        book.setPublisher(param.getPublisher());
        book.setPublishedDate(param.getPublishedDate());
        final List<Author> authors = authorRepository.findAllById(param.getAuthorIds() == null ? List.of() : param.getAuthorIds());
        if (authors.size() != (param.getAuthorIds() == null ? 0 : param.getAuthorIds().size())) {
            throw new NotFoundException("one of authors not found");
        }
        book.setAuthors(new HashSet<>(authors));
        book = bookRepository.save(book);

        return BookDTO.Response.from(book,
                book.getAuthors().stream().map(AuthorDTO.Response::from).toList(),
                book.getReviews().stream().map(ReviewDTO.Response::from).toList());
    }

    public void delete(final Long id) {
        final Book book = bookRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        publisher.publishEvent(new BeforeDeleteBook(id));
        bookRepository.delete(book);
    }

    @EventListener(BeforeDeleteAuthor.class)
    public void on(final BeforeDeleteAuthor event) {
        // remove many-to-many relations at owning side
        bookRepository.findAllByAuthorsId(event.getId()).forEach(book ->
                book.getAuthors().removeIf(author -> author.getId().equals(event.getId())));
    }
}
