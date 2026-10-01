import java.util.*;

class MedianFinder {
    private final PriorityQueue<Integer> low = new PriorityQueue<>(Collections.reverseOrder());
    private final PriorityQueue<Integer> high = new PriorityQueue<>();

    public void add(int x) {
        low.add(x);
        high.add(low.poll());
        if (high.size() > low.size()) {
            low.add(high.poll());
        }
    }

    public double median() {
        if (low.size() > high.size()) {
            return low.peek();
        }
        return (low.peek() + (long) high.peek()) / 2.0;
    }
}
