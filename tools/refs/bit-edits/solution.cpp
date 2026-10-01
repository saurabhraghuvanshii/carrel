int getBit(int n, int i) {
    return n >> i & 1;
}

int setBit(int n, int i) {
    return n | 1 << i;
}

int clearBit(int n, int i) {
    return n & ~(1 << i);
}

int updateBit(int n, int i, int b) {
    return clearBit(n, i) | b << i;
}
