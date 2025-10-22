package com.example.catalog.rest;

import com.example.catalog.model.ProductDTO;
import com.example.catalog.service.CatalogService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.net.URI;
import java.util.List;

@RestController
@RequestMapping(value = "/api/v1/catalogs", produces = MediaType.APPLICATION_JSON_VALUE)
@RequiredArgsConstructor
public class CatalogResource {

    private final CatalogService catalogService;

    @PostMapping("/products")
    public ResponseEntity<ProductDTO.Response> registerProduct(@RequestBody @Valid ProductDTO.RegisterRequest payload) {
        ProductDTO.Response response = catalogService.registerProduct(payload);
        return ResponseEntity.created(URI.create("/api/v1/catalogs/products/" + response.getId())).body(response);
    }

    @GetMapping("/products/{productId}")
    public ResponseEntity<ProductDTO.Response> getProductById(@PathVariable String productId) {
        ProductDTO.Response response = catalogService.findProductById(productId);
        return ResponseEntity.ok(response);
    }

    @GetMapping("/sellers/{sellerId}/products")
    public ResponseEntity<List<ProductDTO.Response>> getProductsBySellerId(@PathVariable String sellerId) {
        List<ProductDTO.Response> products = catalogService.findProductsBySellerId(sellerId);
        return ResponseEntity.ok(products);
    }

    @DeleteMapping("/products/{productId}")
    public ResponseEntity<Void> deleteProduct(@PathVariable String productId) {
        catalogService.deleteProduct(productId);
        return ResponseEntity.noContent().build();
    }

    @PostMapping("/products/{productId}/decreaseStockCount")
    public ResponseEntity<ProductDTO.Response> decreaseStockCount(@PathVariable String productId,
                                                                  @RequestBody @Valid ProductDTO.DecreaseStockRequest payload) {
        ProductDTO.Response response = catalogService.decreaseStockCount(productId, payload.getDecreaseCount());
        return ResponseEntity.ok(response);
    }
}
