package com.example.order.service;

import com.example.order.model.ProductDTO;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

@FeignClient(name = "CATALOG-SERVICE")
public interface CatalogClient {

    @GetMapping("/api/v1/catalogs/products/{productId}")
    ProductDTO.Response getProductById(@PathVariable("productId") String productId);
}
