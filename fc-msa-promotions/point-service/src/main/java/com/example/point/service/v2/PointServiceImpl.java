package com.example.point.service.v2;

import com.example.point.domain.Point;
import com.example.point.model.PointDTO;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

@Service("PointServiceImplV2")
@Slf4j
@RequiredArgsConstructor
public class PointServiceImpl implements PointService {
    @Override
    public Point earn(PointDTO.EarnRequest request) {
        return null;
    }

    @Override
    public Point use(PointDTO.UseRequest request) {
        return null;
    }

    @Override
    public Point cancel(PointDTO.CancelRequest request) {
        return null;
    }

    @Override
    public Long getBalance() {
        return 0L;
    }

    @Override
    public Page<Point> getHistory(Pageable pageable) {
        return null;
    }
}
