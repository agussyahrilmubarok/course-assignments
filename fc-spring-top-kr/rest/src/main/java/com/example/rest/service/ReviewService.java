package com.example.rest.service;

import com.example.rest.domain.Book;
import com.example.rest.domain.Review;
import com.example.rest.events.BeforeDeleteBook;
import com.example.rest.model.ReviewDTO;
import com.example.rest.repos.BookRepository;
import com.example.rest.repos.ReviewRepository;
import com.example.rest.util.NotFoundException;
import com.example.rest.util.ReferencedException;
import org.springframework.context.event.EventListener;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
public class ReviewService {

    private final ReviewRepository reviewRepository;
    private final BookRepository bookRepository;

    public ReviewService(final ReviewRepository reviewRepository,
                         final BookRepository bookRepository) {
        this.reviewRepository = reviewRepository;
        this.bookRepository = bookRepository;
    }

    public List<ReviewDTO.Response> findAll() {
        final List<Review> reviews = reviewRepository.findAll(Sort.by("id"));
        return reviews.stream()
                .map(ReviewDTO.Response::from)
                .toList();
    }

    public ReviewDTO.Response get(final Long id) {
        return reviewRepository.findById(id)
                .map(ReviewDTO.Response::from)
                .orElseThrow(NotFoundException::new);
    }

    public ReviewDTO.Response create(final ReviewDTO.Request param) {
        final Review review = new Review();
        review.setContent(param.getContent());
        review.setRating(param.getRating());
        final Book book = param.getBookId() == null ? null : bookRepository.findById(param.getBookId())
                .orElseThrow(() -> new NotFoundException("book not found"));
        review.setBook(book);
        return ReviewDTO.Response.from(reviewRepository.save(review));
    }

    public ReviewDTO.Response update(final Long id, final ReviewDTO.Request param) {
        final Review review = reviewRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        review.setContent(param.getContent());
        review.setRating(param.getRating());
        final Book book = param.getBookId() == null ? null : bookRepository.findById(param.getBookId())
                .orElseThrow(() -> new NotFoundException("book not found"));
        review.setBook(book);
        return ReviewDTO.Response.from(reviewRepository.save(review));
    }

    public void delete(final Long id) {
        final Review review = reviewRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        reviewRepository.delete(review);
    }

    @EventListener(BeforeDeleteBook.class)
    public void on(final BeforeDeleteBook event) {
        final ReferencedException referencedException = new ReferencedException();
        final Review bookReview = reviewRepository.findFirstByBookId(event.getId());
        if (bookReview != null) {
            referencedException.setKey("book.review.book.referenced");
            referencedException.addParam(bookReview.getId());
            throw referencedException;
        }
    }
}
