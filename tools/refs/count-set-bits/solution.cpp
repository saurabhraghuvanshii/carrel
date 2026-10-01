int countSetBits(int n) {
    unsigned int u = n;
    int count = 0;
    while (u) {
        u &= u - 1;
        count++;
    }
    return count;
}
