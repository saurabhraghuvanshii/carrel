#include <string>
using namespace std;

string longestMirror(string& word) {
    int n = word.size(), best = 0, from = 0;
    for (int c = 0; c < 2 * n - 1; c++) {
        int lo = c / 2, hi = c / 2 + c % 2;
        while (lo >= 0 && hi < n && word[lo] == word[hi]) lo--, hi++;
        if (hi - lo - 1 > best) best = hi - lo - 1, from = lo + 1;
    }
    return word.substr(from, best);
}
