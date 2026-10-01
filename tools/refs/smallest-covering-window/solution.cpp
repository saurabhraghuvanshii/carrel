#include <string>
#include <vector>
using namespace std;

string smallestWindow(string& text, string& pattern) {
    vector<int> need(128);
    for (char ch : pattern) need[ch]++;
    int missing = pattern.size(), bestStart = 0, bestLen = -1, start = 0;
    for (int i = 0; i < (int)text.size(); i++) {
        if (need[text[i]]-- > 0) missing--;
        while (missing == 0) {
            if (bestLen < 0 || i - start + 1 < bestLen) {
                bestStart = start;
                bestLen = i - start + 1;
            }
            if (++need[text[start++]] > 0) missing++;
        }
    }
    return bestLen < 0 ? "" : text.substr(bestStart, bestLen);
}
