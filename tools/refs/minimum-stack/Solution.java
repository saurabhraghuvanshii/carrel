import java.util.*;

class MinStack {
    private final ArrayDeque<Integer> values = new ArrayDeque<>();
    private final ArrayDeque<Integer> mins = new ArrayDeque<>();

    public void push(int x) {
        values.push(x);
        if (mins.isEmpty() || x <= mins.peek()) {
            mins.push(x);
        }
    }

    public void pop() {
        if (values.pop().equals(mins.peek())) {
            mins.pop();
        }
    }

    public int top() {
        return values.peek();
    }

    public int minimum() {
        return mins.peek();
    }
}
