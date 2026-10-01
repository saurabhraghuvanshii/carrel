import java.util.*;

// Ported from the owner's stack/StockSpanqfor.java. Fixed: equal prices now
// count toward the span (>= instead of >).
class Solution {
    public int[] spans(int[] stocks) {
        int[] span = new int[stocks.length];
        Stack<Integer> s = new Stack<>();
        span[0] = 1;
        s.push(0);
        for (int i = 1; i < stocks.length; i++) {
            int currPrice = stocks[i];
            while (!s.isEmpty() && currPrice >= stocks[s.peek()]) {
                s.pop();
            }
            if (s.isEmpty()) {
                span[i] = i + 1;
            } else {
                int prevHigh = s.peek();
                span[i] = i - prevHigh;
            }
            s.push(i);
        }
        return span;
    }
}
