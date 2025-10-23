package com.example.point.service;

import com.example.point.domain.Point;
import com.example.point.model.PointDTO;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

public interface PointService {

    Point earn(PointDTO.EarnRequest request);

    Point use(PointDTO.UseRequest request);

    Point cancel(PointDTO.CancelRequest request);

    Long getBalance();

    Page<Point> getHistory(Pageable pageable);
}
