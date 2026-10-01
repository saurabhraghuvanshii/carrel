import java.util.*;

class Solution {
    public int[] survivors(int[] robots) {
        int[] stack = new int[robots.length];
        int size = 0;
        for (int r : robots) {
            boolean alive = true;
            while (alive && r < 0 && size > 0 && stack[size - 1] > 0) {
                if (stack[size - 1] < -r) {
                    size--;
                } else {
                    if (stack[size - 1] == -r) {
                        size--;
                    }
                    alive = false;
                }
            }
            if (alive) {
                stack[size++] = r;
            }
        }
        return Arrays.copyOf(stack, size);
    }
}
