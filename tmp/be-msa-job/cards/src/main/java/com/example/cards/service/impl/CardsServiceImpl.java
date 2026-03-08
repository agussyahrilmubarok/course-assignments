package com.example.cards.service.impl;

import com.example.cards.constants.CardsConstants;
import com.example.cards.dto.CardsDto;
import com.example.cards.entity.Cards;
import com.example.cards.exception.CardAlreadyExistsException;
import com.example.cards.exception.ResourceNotFoundException;
import com.example.cards.mapper.CardsMapper;
import com.example.cards.repository.CardsRepository;
import com.example.cards.service.ICardsService;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.Optional;
import java.util.Random;

@Service
@Slf4j
@AllArgsConstructor
public class CardsServiceImpl implements ICardsService {

    private CardsRepository cardsRepository;

    /**
     * @param mobileNumber - Mobile Number of the Customer
     */
    @Override
    public void createCard(String mobileNumber) {
        log.info("Attempting to create a new card for mobileNumber: {}", mobileNumber);
        Optional<Cards> optionalCards = cardsRepository.findByMobileNumber(mobileNumber);
        if (optionalCards.isPresent()) {
            log.error("Card already exists for mobileNumber: {}", mobileNumber);
            throw new CardAlreadyExistsException("Card already registered with given mobileNumber " + mobileNumber);
        }

        Cards newCard = createNewCard(mobileNumber);
        cardsRepository.save(newCard);
        log.info("Card successfully created with cardNumber: {} for mobileNumber: {}", newCard.getCardNumber(), mobileNumber);
    }

    /**
     * @param mobileNumber - Mobile Number of the Customer
     * @return the new card details
     */
    private Cards createNewCard(String mobileNumber) {
        long randomCardNumber = 100000000000L + new Random().nextInt(900000000);
        Cards newCard = new Cards();
        newCard.setCardNumber(Long.toString(randomCardNumber));
        newCard.setMobileNumber(mobileNumber);
        newCard.setCardType(CardsConstants.CREDIT_CARD);
        newCard.setTotalLimit(CardsConstants.NEW_CARD_LIMIT);
        newCard.setAmountUsed(0);
        newCard.setAvailableAmount(CardsConstants.NEW_CARD_LIMIT);
        log.debug("New card initialized for mobileNumber: {}, cardNumber: {}", mobileNumber, randomCardNumber);
        return newCard;
    }

    /**
     * @param mobileNumber - Input mobile Number
     * @return Card Details based on a given mobileNumber
     */
    @Override
    public CardsDto fetchCard(String mobileNumber) {
        log.info("Fetching card for mobileNumber: {}", mobileNumber);
        Cards cards = cardsRepository.findByMobileNumber(mobileNumber).orElseThrow(() -> {
            log.error("Card not found for mobileNumber: {}", mobileNumber);
            return new ResourceNotFoundException("Card", "mobileNumber", mobileNumber);
        });

        CardsDto cardsDto = CardsMapper.mapToCardsDto(cards, new CardsDto());
        log.info("Successfully fetched card details for mobileNumber: {}", mobileNumber);
        return cardsDto;
    }

    /**
     * @param cardsDto - CardsDto Object
     * @return boolean indicating if the update of card details is successful or not
     */
    @Override
    public boolean updateCard(CardsDto cardsDto) {
        log.info("Updating card with cardNumber: {}", cardsDto.getCardNumber());
        Cards cards = cardsRepository.findByCardNumber(cardsDto.getCardNumber()).orElseThrow(() -> {
            log.error("Card not found for cardNumber: {}", cardsDto.getCardNumber());
            return new ResourceNotFoundException("Card", "CardNumber", cardsDto.getCardNumber());
        });

        CardsMapper.mapToCards(cardsDto, cards);
        cardsRepository.save(cards);
        log.info("Card updated successfully for cardNumber: {}", cardsDto.getCardNumber());
        return true;
    }

    /**
     * @param mobileNumber - Input MobileNumber
     * @return boolean indicating if the delete of card details is successful or not
     */
    @Override
    public boolean deleteCard(String mobileNumber) {
        log.info("Deleting card for mobileNumber: {}", mobileNumber);
        Cards cards = cardsRepository.findByMobileNumber(mobileNumber).orElseThrow(() -> {
            log.error("Card not found for mobileNumber: {}", mobileNumber);
            return new ResourceNotFoundException("Card", "mobileNumber", mobileNumber);
        });

        cardsRepository.deleteById(cards.getCardId());
        log.info("Card deleted successfully for mobileNumber: {}", mobileNumber);
        return true;
    }
}
