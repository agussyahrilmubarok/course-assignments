package com.example.grpc.client.rest;

import com.example.grpc.client.model.ReviewDTO;
import com.example.grpc.client.service.ReviewService;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import jakarta.validation.Valid;
import org.springframework.hateoas.CollectionModel;
import org.springframework.hateoas.EntityModel;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

import static org.springframework.hateoas.server.mvc.WebMvcLinkBuilder.linkTo;
import static org.springframework.hateoas.server.mvc.WebMvcLinkBuilder.methodOn;


@RestController
@RequestMapping(value = "/api/reviews", produces = MediaType.APPLICATION_JSON_VALUE)
public class ReviewResource {

    private final ReviewService reviewService;

    public ReviewResource(final ReviewService reviewService) {
        this.reviewService = reviewService;
    }

    @GetMapping
    public ResponseEntity<CollectionModel<EntityModel<ReviewDTO.Response>>> getAllReviews() {
        final List<EntityModel<ReviewDTO.Response>> reviews = reviewService.findAll()
                .stream()
                .map(review -> EntityModel.of(review,
                        linkTo(methodOn(this.getClass()).getReview(review.getId())).withSelfRel(),
                        linkTo(methodOn(this.getClass()).getAllReviews()).withRel("all-reviews")))
                .toList();

        final CollectionModel<EntityModel<ReviewDTO.Response>> collectionModel = CollectionModel.of(reviews,
                linkTo(methodOn(this.getClass()).getAllReviews()).withSelfRel());

        return ResponseEntity.ok(collectionModel);
    }

    @GetMapping("/{id}")
    public ResponseEntity<EntityModel<ReviewDTO.Response>> getReview(@PathVariable(name = "id") final Long id) {
        final ReviewDTO.Response reviewDTO = reviewService.get(id);

        final EntityModel<ReviewDTO.Response> entityModel = EntityModel.of(reviewDTO);
        entityModel.add(linkTo(methodOn(this.getClass()).getReview(id)).withSelfRel());
        entityModel.add(linkTo(methodOn(this.getClass()).getAllReviews()).withRel("all-reviews"));

        return ResponseEntity.ok(entityModel);
    }

    @PostMapping
    @ApiResponse(responseCode = "201")
    public ResponseEntity<EntityModel<ReviewDTO.Response>> createReview(@RequestBody @Valid final ReviewDTO.Request payload) {
        final ReviewDTO.Response reviewDTO = reviewService.create(payload);

        final EntityModel<ReviewDTO.Response> entityModel = EntityModel.of(reviewDTO);
        entityModel.add(linkTo(methodOn(this.getClass()).getAllReviews()).withRel("all-reviews"));
        entityModel.add(linkTo(methodOn(this.getClass()).getReview(reviewDTO.getId())).withRel("review-by-id"));

        return new ResponseEntity<>(entityModel, HttpStatus.CREATED);
    }

    @PutMapping("/{id}")
    public ResponseEntity<EntityModel<ReviewDTO.Response>> updateReview(@PathVariable(name = "id") final Long id,
                                                                        @RequestBody @Valid final ReviewDTO.Request payload) {
        final ReviewDTO.Response reviewDTO = reviewService.update(id, payload);

        final EntityModel<ReviewDTO.Response> entityModel = EntityModel.of(reviewDTO,
                linkTo(methodOn(this.getClass()).getReview(id)).withSelfRel(),
                linkTo(methodOn(this.getClass()).getAllReviews()).withRel("all-reviews"));

        return ResponseEntity.ok(entityModel);
    }

    @DeleteMapping("/{id}")
    @ApiResponse(responseCode = "204")
    public ResponseEntity<Void> deleteReview(@PathVariable(name = "id") final Long id) {
        reviewService.delete(id);

        return ResponseEntity.noContent().build();
    }
}
