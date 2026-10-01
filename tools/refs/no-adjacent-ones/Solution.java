import java.util.*;

// Ported from the owner's recursion/BinaryString.java (printBinString); it now
// collects the strings instead of printing them.
class Solution {
    private final List<String> out = new ArrayList<>();

    public List<String> noAdjacentOnes(int n) {
        printBinString(n, 0, "");
        return out;
    }

    private void printBinString(int n, int lastPlace, String str) {
        if (n == 0) {
            out.add(str);
            return;
        }
        printBinString(n - 1, 0, str + "0");
        if (lastPlace == 0) {
            printBinString(n - 1, 1, str + "1");
        }
    }
}
