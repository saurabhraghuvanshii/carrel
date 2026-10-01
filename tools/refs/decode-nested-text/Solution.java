import java.util.*;

class Solution {
    public String decode(String text) {
        ArrayDeque<Integer> counts = new ArrayDeque<>();
        ArrayDeque<StringBuilder> outer = new ArrayDeque<>();
        StringBuilder current = new StringBuilder();
        int k = 0;
        for (char ch : text.toCharArray()) {
            if (Character.isDigit(ch)) {
                k = k * 10 + (ch - '0');
            } else if (ch == '[') {
                counts.push(k);
                outer.push(current);
                current = new StringBuilder();
                k = 0;
            } else if (ch == ']') {
                String inner = current.toString();
                current = outer.pop();
                current.append(inner.repeat(counts.pop()));
            } else {
                current.append(ch);
            }
        }
        return current.toString();
    }
}
