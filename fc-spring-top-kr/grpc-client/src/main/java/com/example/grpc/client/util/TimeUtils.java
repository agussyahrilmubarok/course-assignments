package com.example.grpc.client.util;

import com.google.protobuf.Timestamp;

import java.time.Instant;
import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.time.ZoneId;

public final class TimeUtils {

    private TimeUtils() {
    }

    public static OffsetDateTime tsToOffsetDateTime(Timestamp ts) {
        if (ts == null) {
            return null;
        }
        Instant instant = Instant.ofEpochSecond(ts.getSeconds(), ts.getNanos());
        ZoneId zoneId = ZoneId.of("Asia/Jakarta");
        return instant.atZone(zoneId).toOffsetDateTime();
    }

    public static LocalDate tsToLocalDate(Timestamp ts) {
        if (ts == null) {
            return null;
        }
        Instant instant = Instant.ofEpochSecond(ts.getSeconds(), ts.getNanos());
        ZoneId zoneId = ZoneId.of("Asia/Jakarta");
        return instant.atZone(zoneId).toLocalDate();
    }

    public static Timestamp toTimestamp(OffsetDateTime odt) {
        if (odt == null) {
            return null;
        }
        Instant instant = odt.toInstant();
        return Timestamp.newBuilder()
                .setSeconds(instant.getEpochSecond())
                .setNanos(instant.getNano())
                .build();
    }

    public static Timestamp toTimestamp(LocalDate date) {
        if (date == null) {
            return null;
        }
        Instant instant = date.atStartOfDay(ZoneId.of("Asia/Jakarta")).toInstant();
        return Timestamp.newBuilder()
                .setSeconds(instant.getEpochSecond())
                .setNanos(instant.getNano())
                .build();
    }
}
