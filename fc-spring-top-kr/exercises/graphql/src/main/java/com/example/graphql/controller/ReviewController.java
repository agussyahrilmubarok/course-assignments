package com.example.graphql.controller;

import com.example.graphql.model.ReviewDTO;
import com.example.graphql.service.ReviewService;
import lombok.RequiredArgsConstructor;
import org.springframework.graphql.data.method.annotation.*;
import org.springframework.stereotype.Controller;
import java.util.List;

@Controller
@RequiredArgsConstructor
public class ReviewController {

    private final ReviewService reviewService;

    @QueryMapping
    public List<ReviewDTO.Response> findAllReviews() {
        return reviewService.findAll();
    }

    @QueryMapping
    public ReviewDTO.Response getReview(@Argument Long id) {
        return reviewService.get(id);
    }

    @MutationMapping
    public ReviewDTO.Response createReview(@Argument("input") ReviewDTO.ReviewRequest input) {
        return reviewService.create(input);
    }

    @MutationMapping
    public ReviewDTO.Response updateReview(@Argument Long id, @Argument("input") ReviewDTO.ReviewRequest input) {
        return reviewService.update(id, input);
    }

    @MutationMapping
    public Boolean deleteReview(@Argument Long id) {
        reviewService.delete(id);
        return true;
    }
}
