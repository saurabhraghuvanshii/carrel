class Solution {
    public int textToNumber(String text) {
        int i = 0, n = text.length();
        while (i < n && text.charAt(i) == ' ') {
            i++;
        }
        int sign = 1;
        if (i < n && (text.charAt(i) == '+' || text.charAt(i) == '-')) {
            sign = text.charAt(i) == '-' ? -1 : 1;
            i++;
        }
        long value = 0;
        while (i < n && Character.isDigit(text.charAt(i))) {
            value = value * 10 + (text.charAt(i) - '0');
            if (value > Integer.MAX_VALUE) {
                return sign == 1 ? Integer.MAX_VALUE : Integer.MIN_VALUE;
            }
            i++;
        }
        return (int) (sign * value);
    }
}
