package com.example.rest.rest;

import com.example.rest.model.AuthorDTO;
import com.example.rest.service.AuthorService;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import jakarta.validation.Valid;
import org.springframework.hateoas.CollectionModel;
import org.springframework.hateoas.EntityModel;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.stream.Collectors;

import static org.springframework.hateoas.server.mvc.WebMvcLinkBuilder.linkTo;
import static org.springframework.hateoas.server.mvc.WebMvcLinkBuilder.methodOn;

@RestController
@RequestMapping(value = "/api/authors", produces = MediaType.APPLICATION_JSON_VALUE)
public class AuthorResource {

    private final AuthorService authorService;

    public AuthorResource(final AuthorService authorService) {
        this.authorService = authorService;
    }

    @GetMapping
    public ResponseEntity<CollectionModel<EntityModel<AuthorDTO.Response>>> getAllAuthors() {
        final List<EntityModel<AuthorDTO.Response>> authors = authorService.findAll()
                .stream()
                .map(author -> EntityModel.of(author,
                        linkTo(methodOn(this.getClass()).getAuthor(author.getId())).withSelfRel(),
                        linkTo(methodOn(this.getClass()).getAllAuthors()).withRel("all-authors")))
                .collect(Collectors.toList());
        final CollectionModel<EntityModel<AuthorDTO.Response>> collectionModel = CollectionModel.of(authors,
                linkTo(methodOn(this.getClass()).getAllAuthors()).withSelfRel());

        return ResponseEntity.ok(collectionModel);
    }

    @GetMapping("/{id}")
    public ResponseEntity<EntityModel<AuthorDTO.Response>> getAuthor(@PathVariable(name = "id") final Long id) {
        final AuthorDTO.Response authorDTO = authorService.get(id);
        final EntityModel<AuthorDTO.Response> entityModel = EntityModel.of(authorDTO);
        entityModel.add(linkTo(methodOn(this.getClass()).getAuthor(id)).withSelfRel());
        entityModel.add(linkTo(methodOn(this.getClass()).getAllAuthors()).withRel("all-authors"));

        return ResponseEntity.ok(entityModel);
    }

    @PostMapping
    @ApiResponse(responseCode = "201")
    public ResponseEntity<EntityModel<AuthorDTO.Response>> createAuthor(@RequestBody @Valid final AuthorDTO.Request payload) {
        final AuthorDTO.Response authorDTO = authorService.create(payload);
        final EntityModel<AuthorDTO.Response> entityModel = EntityModel.of(authorDTO);
        entityModel.add(linkTo(methodOn(this.getClass()).getAllAuthors()).withRel("all-authors"));
        entityModel.add(linkTo(methodOn(this.getClass()).getAuthor(authorDTO.getId())).withRel("author-by-id"));

        return new ResponseEntity<>(entityModel, HttpStatus.CREATED);
    }

    @PutMapping("/{id}")
    public ResponseEntity<EntityModel<AuthorDTO.Response>> updateAuthor(@PathVariable(name = "id") final Long id,
                                                                        @RequestBody @Valid final AuthorDTO.Request payload) {
        final AuthorDTO.Response authorDTO = authorService.update(id, payload);
        final EntityModel<AuthorDTO.Response> entityModel = EntityModel.of(authorDTO,
                linkTo(methodOn(this.getClass()).getAuthor(id)).withSelfRel(),
                linkTo(methodOn(this.getClass()).getAllAuthors()).withRel("all-authors"));

        return ResponseEntity.ok(entityModel);
    }

    @DeleteMapping("/{id}")
    @ApiResponse(responseCode = "204")
    public ResponseEntity<Void> deleteAuthor(@PathVariable(name = "id") final Long id) {
        authorService.delete(id);

        return ResponseEntity.noContent().build();
    }
}
