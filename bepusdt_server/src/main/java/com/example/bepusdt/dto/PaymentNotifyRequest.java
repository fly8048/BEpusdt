package com.example.bepusdt.dto;

import java.math.BigDecimal;

public class PaymentNotifyRequest {

    private String trade_id;
    private String order_id;
    private BigDecimal amount;
    private BigDecimal actual_amount;
    private String token;
    private String block_transaction_id;
    private Integer status;
    private String signature;

    // ===== getter / setter =====

    public String getTrade_id() {
        return trade_id;
    }

    public void setTrade_id(String trade_id) {
        this.trade_id = trade_id;
    }

    public String getOrder_id() {
        return order_id;
    }

    public void setOrder_id(String order_id) {
        this.order_id = order_id;
    }

    public BigDecimal getAmount() {
        return amount;
    }

    public void setAmount(BigDecimal amount) {
        this.amount = amount;
    }

    public BigDecimal getActual_amount() {
        return actual_amount;
    }

    public void setActual_amount(BigDecimal actual_amount) {
        this.actual_amount = actual_amount;
    }

    public String getToken() {
        return token;
    }

    public void setToken(String token) {
        this.token = token;
    }

    public String getBlock_transaction_id() {
        return block_transaction_id;
    }

    public void setBlock_transaction_id(String block_transaction_id) {
        this.block_transaction_id = block_transaction_id;
    }

    public Integer getStatus() {
        return status;
    }

    public void setStatus(Integer status) {
        this.status = status;
    }

    public String getSignature() {
        return signature;
    }

    public void setSignature(String signature) {
        this.signature = signature;
    }
}

