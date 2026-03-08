package com.finly.finly.util;

import org.bson.Document;
import org.springframework.core.convert.converter.Converter;

import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.Date;

import static com.finly.finly.util.MongoOffsetDateTimeWriter.DATE_FIELD;
import static com.finly.finly.util.MongoOffsetDateTimeWriter.OFFSET_FIELD;


public class MongoOffsetDateTimeReader implements Converter<Document, OffsetDateTime> {

    @Override
    public OffsetDateTime convert(final Document document) {
        final Date dateTime = document.getDate(DATE_FIELD);
        final ZoneOffset offset = ZoneOffset.of(document.getString(OFFSET_FIELD));
        return OffsetDateTime.ofInstant(dateTime.toInstant(), offset);
    }

}
