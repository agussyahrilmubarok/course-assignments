package com.example.graphql.service;

import com.example.graphql.domain.Book;
import com.example.graphql.domain.Review;
import com.example.graphql.events.BeforeDeleteBook;
import com.example.graphql.model.ReviewDTO;
import com.example.graphql.repos.BookRepository;
import com.example.graphql.repos.ReviewRepository;
import com.example.graphql.util.NotFoundException;
import com.example.graphql.util.ReferencedException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.event.EventListener;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@Slf4j
public class ReviewService {

    private final ReviewRepository reviewRepository;
    private final BookRepository bookRepository;

    public ReviewService(final ReviewRepository reviewRepository,
                         final BookRepository bookRepository) {
        this.reviewRepository = reviewRepository;
        this.bookRepository = bookRepository;
    }

    public List<ReviewDTO.Response> findAll() {
        log.info("Fetching all reviews");
        final List<Review> reviews = reviewRepository.findAll(Sort.by("id"));
        log.debug("Found {} reviews", reviews.size());
        return reviews.stream()
                .map(ReviewDTO.Response::from)
                .toList();
    }

    public ReviewDTO.Response get(final Long id) {
        log.info("Fetching review with id: {}", id);
        return reviewRepository.findById(id)
                .map(review -> {
                    log.debug("Review found with content: {}", review.getContent());
                    return ReviewDTO.Response.from(review);
                })
                .orElseThrow(() -> {
                    log.error("Review with id {} not found", id);
                    return new NotFoundException();
                });
    }

    public ReviewDTO.Response create(final ReviewDTO.ReviewRequest param) {
        log.info("Creating new review for bookId: {}", param.getBookId());
        final Review review = new Review();
        review.setContent(param.getContent());
        review.setRating(param.getRating());

        final Book book = param.getBookId() == null ? null : bookRepository.findById(param.getBookId())
                .orElseThrow(() -> {
                    log.error("Book with id {} not found for review creation", param.getBookId());
                    return new NotFoundException("book not found");
                });
        review.setBook(book);

        final Review saved = reviewRepository.save(review);
        log.debug("Review created with id: {}", saved.getId());
        return ReviewDTO.Response.from(saved);
    }

    public ReviewDTO.Response update(final Long id, final ReviewDTO.ReviewRequest param) {
        log.info("Updating review with id: {}", id);
        final Review review = reviewRepository.findById(id)
                .orElseThrow(() -> {
                    log.error("Review with id {} not found for update", id);
                    return new NotFoundException();
                });

        review.setContent(param.getContent());
        review.setRating(param.getRating());

        final Book book = param.getBookId() == null ? null : bookRepository.findById(param.getBookId())
                .orElseThrow(() -> {
                    log.error("Book with id {} not found for review update", param.getBookId());
                    return new NotFoundException("book not found");
                });
        review.setBook(book);

        final Review updated = reviewRepository.save(review);
        log.debug("Review with id {} updated successfully", id);
        return ReviewDTO.Response.from(updated);
    }

    public void delete(final Long id) {
        log.info("Deleting review with id: {}", id);
        final Review review = reviewRepository.findById(id)
                .orElseThrow(() -> {
                    log.error("Review with id {} not found for deletion", id);
                    return new NotFoundException();
                });
        reviewRepository.delete(review);
        log.debug("Review with id {} deleted successfully", id);
    }

    @EventListener(BeforeDeleteBook.class)
    public void on(final BeforeDeleteBook event) {
        log.info("Handling BeforeDeleteBook event for bookId: {}", event.getId());
        final Review bookReview = reviewRepository.findFirstByBookId(event.getId());
        if (bookReview != null) {
            log.warn("Cannot delete book {} because it has a review with id {}", event.getId(), bookReview.getId());
            final ReferencedException referencedException = new ReferencedException();
            referencedException.setKey("book.review.book.referenced");
            referencedException.addParam(bookReview.getId());
            throw referencedException;
        }
        log.debug("No reviews found for bookId: {}, deletion allowed", event.getId());
    }
}
