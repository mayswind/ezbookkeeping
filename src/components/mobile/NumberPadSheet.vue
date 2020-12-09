<template>
    <f7-sheet class="numpad-sheet" :opened="show" @sheet:closed="onSheetClosed">
        <f7-page-content class="no-margin no-padding">
            <f7-row class="numpad-values">
                <span class="numpad-value">{{ currentDisplay }}</span>
            </f7-row>
            <f7-row class="numpad-buttons">
                <span class="numpad-button numpad-button-num" @click="inputNum(7)">
                    <span class="numpad-button-text">7</span>
                </span>
                <span class="numpad-button numpad-button-num" @click="inputNum(8)">
                    <span class="numpad-button-text">8</span>
                </span>
                <span class="numpad-button numpad-button-num" @click="inputNum(9)">
                    <span class="numpad-button-text">9</span>
                </span>
                <span class="numpad-button numpad-button-function no-right-border" @click="setSymbol('×')">
                    <span class="numpad-button-text">&times;</span>
                </span>
                <span class="numpad-button numpad-button-num" @click="inputNum(4)">
                    <span class="numpad-button-text">4</span>
                </span>
                <span class="numpad-button numpad-button-num" @click="inputNum(5)">
                    <span class="numpad-button-text">5</span>
                </span>
                <span class="numpad-button numpad-button-num" @click="inputNum(6)">
                    <span class="numpad-button-text">6</span>
                </span>
                <span class="numpad-button numpad-button-function no-right-border" @click="setSymbol('−')">
                    <span class="numpad-button-text">&minus;</span>
                </span>
                <span class="numpad-button numpad-button-num" @click="inputNum(1)">
                    <span class="numpad-button-text">1</span>
                </span>
                <span class="numpad-button numpad-button-num" @click="inputNum(2)">
                    <span class="numpad-button-text">2</span>
                </span>
                <span class="numpad-button numpad-button-num" @click="inputNum(3)">
                    <span class="numpad-button-text">3</span>
                </span>
                <span class="numpad-button numpad-button-function no-right-border" @click="setSymbol('+')">
                    <span class="numpad-button-text">&plus;</span>
                </span>
                <span class="numpad-button numpad-button-num no-bottom-border" @click="inputDot()">
                    <span class="numpad-button-text">.</span>
                </span>
                <span class="numpad-button numpad-button-num no-bottom-border" @click="inputNum(0)">
                    <span class="numpad-button-text">0</span>
                </span>
                <span class="numpad-button numpad-button-num no-bottom-border" @click="backspace">
                    <span class="numpad-button-text">
                        <f7-icon f7="delete_left"></f7-icon>
                    </span>
                </span>
                <span class="numpad-button numpad-button-function numpad-button-confirm no-right-border no-bottom-border" @click="confirm()">
                    <span :class="{ 'numpad-button-text': true, 'numpad-button-text-confirm': !currentSymbol }">{{ confirmText }}</span>
                </span>
            </f7-row>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
export default {
    props: [
        'amount',
        'show'
    ],
    data() {
        const self = this;

        return {
            previousValue: '',
            currentSymbol: '',
            currentValue: self.getStringValue(self.amount)
        }
    },
    computed: {
        currentDisplay() {
            const previousValue = this.$utilities.appendThousandsSeparator(this.previousValue);
            const currentValue = this.$utilities.appendThousandsSeparator(this.currentValue);

            if (this.currentSymbol) {
                return `${previousValue} ${this.currentSymbol} ${currentValue}`;
            } else {
                return currentValue;
            }
        },
        confirmText() {
            if (this.currentSymbol) {
                return '=';
            } else {
                return this.$i18n.t('OK');
            }
        }
    },
    watch: {
        amount: function(newValue){
            this.currentValue = this.getStringValue(newValue);
        }
    },
    methods: {
        getStringValue(value) {
            let str = this.$utilities.numericCurrencyToString(value);

            if (str.indexOf(',')) {
                str = str.replaceAll(/,/g, '');
            }

            const dotPos = str.indexOf('.');

            if (dotPos < 0) {
                if (str === '0') {
                    return '';
                }

                return str;
            }

            let integer = str.substr(0, dotPos);
            let decimals = str.substring(dotPos + 1, str.length);
            let newDecimals = '';

            for (let i = decimals.length - 1; i >= 0; i--) {
                if (decimals[i] !== '0' || newDecimals.length > 0) {
                    newDecimals = decimals[i] + newDecimals;
                }
            }

            if (newDecimals.length < 1) {
                if (integer === '0') {
                    return '';
                }

                return integer;
            }

            return `${integer}.${newDecimals}`;
        },
        inputNum(num) {
            if (this.currentValue === '0') {
                this.currentValue = num.toString();
                return;
            } else if (this.currentValue === '-0') {
                this.currentValue = '-' + num.toString();
                return;
            }

            const dotPos = this.currentValue.indexOf('.');

            if (dotPos >= 0 && this.currentValue.substring(dotPos + 1, this.currentValue.length).length >= 2) {
                return;
            }

            this.currentValue = this.currentValue + num.toString();
        },
        inputDot() {
            if (this.currentValue.indexOf('.') >= 0) {
                return;
            }

            if (this.currentValue.length < 1) {
                this.currentValue = '0';
            } else if (this.currentValue === '-') {
                this.currentValue = '-0';
            }

            this.currentValue = this.currentValue + '.';
        },
        setSymbol(symbol) {
            if (this.currentValue) {
                if (this.currentSymbol) {
                    this.confirm();
                }

                this.previousValue = this.currentValue;
                this.currentValue = '';
            }

            this.currentSymbol = symbol;
        },
        backspace() {
            if (!this.currentValue || this.currentValue.length < 1) {
                if (this.currentSymbol) {
                    this.currentValue = this.previousValue;
                    this.previousValue = '';
                    this.currentSymbol = '';
                }

                return;
            }

            this.currentValue = this.currentValue.substr(0, this.currentValue.length - 1);
        },
        confirm() {
            if (this.currentSymbol && this.currentValue.length >= 1) {
                const previousValue = this.$utilities.stringCurrencyToNumeric(this.previousValue);
                const currentValue = this.$utilities.stringCurrencyToNumeric(this.currentValue);
                let finalValue = 0;

                switch (this.currentSymbol) {
                    case '+':
                        finalValue = previousValue + currentValue;
                        break;
                    case '−':
                        finalValue = previousValue - currentValue;
                        break;
                    case '×':
                        finalValue = Math.round(previousValue * currentValue / 100);
                        break;
                    default:
                        finalValue = previousValue;
                }

                this.currentValue = this.getStringValue(finalValue);
                this.previousValue = '';
                this.currentSymbol = '';
            } else if (this.currentSymbol && this.currentValue.length < 1) {
                this.currentValue = this.previousValue;
                this.previousValue = '';
                this.currentSymbol = '';
            } else {
                const amount = this.$utilities.stringCurrencyToNumeric(this.currentValue);
                this.currentValue = '';

                this.$emit('numpad:change', amount);
            }
        },
        onSheetClosed() {
            this.currentValue = '';
            this.previousValue = '';
            this.currentSymbol = '';

            this.$emit('numpad:closed');
        }
    }
}
</script>

<style>
.numpad-sheet {
    height: auto;
}

.numpad-values {
    border-bottom: 1px solid var(--f7-page-bg-color);
}

.numpad-value {
    display: flex;
    position: relative;
    padding-left: 16px;
    font-size: 24px;
    line-height: 1;
    height: 50px;
    justify-content: center;
    align-items: center;
    box-sizing: border-box;
    user-select: none;
}

.numpad-buttons {
    display: flex;
    flex-wrap: wrap;
}

.numpad-button {
    display: flex;
    position: relative;
    text-align: center;
    border-right: 1px solid var(--f7-page-bg-color);
    border-bottom: 1px solid var(--f7-page-bg-color);
    cursor: pointer;
    height: 60px;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    box-sizing: border-box;
    user-select: none;
}

.numpad-button.active-state {
    background-color: var(--f7-list-button-pressed-bg-color);
}

.numpad-button-num {
    width: calc(80% / 3);
}

.numpad-button-function {
    width: 20%;
}

.numpad-button-confirm {
    background-color: var(--f7-theme-color);
    color: #ffffff;
}

.numpad-button-confirm.active-state {
    background-color: var(--f7-theme-color-tint);
}

.numpad-button-text {
    display: block;
    font-size: 28px;
    line-height: 1;
}

.numpad-button-text-confirm {
    font-size: 20px;
}
</style>
