import java.util.*;

// Ported from the owner's Greedy/MinAbsoluteDif.java. Fixed: the sum was an
// int, which overflows past about 2.1 × 10⁹; it is now a long.
class Solution {
    public long smallestGapSum(int[] A, int[] B) {
        Arrays.sort(A);
        Arrays.sort(B);
        long minDiff = 0;
        for (int i = 0; i < A.length; i++) {
            minDiff += Math.abs(A[i] - B[i]);
        }
        return minDiff;
    }
}
