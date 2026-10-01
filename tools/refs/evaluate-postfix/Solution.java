import java.util.*;

class Solution {
    public int evaluate(String[] tokens) {
        ArrayDeque<Integer> stack = new ArrayDeque<>();
        for (String tok : tokens) {
            if (tok.length() == 1 && "+-*/".indexOf(tok.charAt(0)) >= 0) {
                int b = stack.pop();
                int a = stack.pop();
                switch (tok.charAt(0)) {
                    case '+': stack.push(a + b); break;
                    case '-': stack.push(a - b); break;
                    case '*': stack.push(a * b); break;
                    default: stack.push(a / b);
                }
            } else {
                stack.push(Integer.parseInt(tok));
            }
        }
        return stack.pop();
    }
}
