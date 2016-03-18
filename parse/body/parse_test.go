package parse

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var testHTML string = `<div class="asset-double-wide double-wide p402_premium" role="main" itemprop="articleBody"><div id="module-position-OlUUBRQhWAc" class="story-asset story-metadata-asset"><div class="article-metadata-wrap"><section id="module-position-OlUUBRQOYUs" class="storymetadata-bucket expandable-photo-module story-expandable-photo-module"><aside itemprop="associatedMedia" itemscope="" itemtype="http://schema.org/ImageObject" class="single-photo expandable-collapsed"><div class="image-wrap"><img class="expand-img-horiz" itemprop="url" src="http://www.gannett-cdn.com/-mm-/d622cfe6e7e2b406a108c56065a889106ca80f0a/c=0-26-800-626&amp;r=x404&amp;c=534x401/local/-/media/2015/11/17/Livonia/B9319685824Z.1_20151117014940_000_G4NCJ9LKJ.1-0.jpg" alt="JS15no11" data-mycapture-src="http://www.gannett-cdn.com/media/2015/11/17/Livonia/B9319685824Z.1_20151117014940_000_G4NCJ9LKJ.1-0.jpg" data-mycapture-sm-src="http://www.gannett-cdn.com/-mm-/2a6a79749746b60a2de9e79be26ac4ec94dc84cf/r=490x400/local/-/media/2015/11/17/Livonia/B9319685824Z.1_20151117014940_000_G4NCJ9LKJ.1-0.jpg"><span class="toggle"></span><meta itemprop="name" content="JS15no11"><span class="mycapture-small-btn mycapture-btn-with-text mycapture-expandable-photo-btn-small js-mycapture-btn-small">Buy Photo</span></div><p class="image-credit-wrap"><span class="js-caption-wrapper"><span class="cutline js-caption">Senior Jack Cikra is the leading scorer among returning players for the Farmington Hills Unified hockey team.</span><meta itemprop="copyrightHolder" content=""><span class="credit">(Photo: john stormzand | staff photographer)</span></span><span class="mycapture-btn mycapture-btn-with-text mycapture-expandable-photo-btn-large js-mycapture-btn">Buy Photo</span></p></aside></section></div></div><div id="module-position-OlUUBRRO510" class="story-asset inline-share-tools-asset"><div class="inline-share-tools asset-inline-share-tools asset-inline-share-tools-top"><span class="inline-share-btn inline-share-btn-facebook" data-share-method="facebook" data-uotrack="InlineShareFacebookBtn"><span class="inline-share-count inline-share-count-facebook"></span><span class="inline-share-btn-label inline-share-btn-label-facebook">CONNECT</span></span><a class="inline-share-btn inline-share-btn-twitter" data-share-method="twitter" data-uotrack="InlineShareTwitterLink" data-popup-width="550" data-popup-height="450" href="https://twitter.com/intent/tweet?url=http%3A//www.hometownlife.com/story/sports/high-school/farmington/2015/11/17/flyers-need-find-scoring-punch/75912786/&amp;text=Flyers%20need%20to%20find%20more%20scoring%20punch&amp;via=hometownlife" target="_blank" data-popup=""><span class="inline-share-count inline-share-count-twitter"> 4 </span><span class="inline-share-btn-label inline-share-btn-label-twitter">TWEET</span></a><a class="inline-share-btn inline-share-btn-linkedin" data-share-method="linkedin" data-popup-width="600" data-popup-height="455" data-uotrack="UtilityBarFlyoutLinkedInLink" href="http://www.linkedin.com/shareArticle?url=http%3A//www.hometownlife.com/story/sports/high-school/farmington/2015/11/17/flyers-need-find-scoring-punch/75912786/&amp;mini=true" target="_blank" data-popup=""><span class="inline-share-count inline-share-count-linkedin"></span><span class="inline-share-btn-label inline-share-btn-label-linkedin">LINKEDIN</span></a><span class="inline-share-btn inline-share-btn-comments" data-share-method="comments" data-uotrack="InlineShareCommentsBtn"><span class="inline-share-count inline-share-count-comments"></span><span class="inline-share-btn-label inline-share-btn-label-comments">COMMENT</span></span><span class="inline-share-btn inline-share-btn-email" data-share-method="email" data-uotrack="InlineShareEmailBtn"><span class="inline-share-count inline-share-count-email"></span><span class="inline-share-btn-label inline-share-btn-label-email">EMAIL</span></span><span class="inline-share-btn inline-share-btn-more" data-share-method="facebook" data-uotrack="InlineShareMoreBtn"><span class="inline-share-count inline-share-count-more"></span><span class="inline-share-btn-label inline-share-btn-label-more">MORE</span></span></div></div><p>On the eve of a new hockey season, the major issue facing Farmington Hills Unified coach Ken Anderson and his team is how to replace all the scoring that was lost.</p><p>Austin Bottrell was an all-state player after the greatest offensive output by a Flyers player. He pumped in 34 goals and had 28 assists for 62 points.</p><p>Danny Arnold and his 35 points was another graduation loss, and Joey Lajcaj (23 points) opted to play junior hockey in Colorado instead of returning for his senior season.</p><p>“It’s going to be a struggle to make up that many points,” Anderson said. “It will have to be a team effort to get some of them. We’ll have to play a defensive style game and take care of the defensive zone first.”</p><p>First of all, the Flyers will look to senior Jack Cikra, the leading scorer among the returning players. He had 17 goals and 24 points to rank third on the team.</p><p>Seniors Brandon Glasser (15) and Andrew Nathan (18) will play prominent roles in helping to recharge the offense, too. Glasser had five goals and Nathan six.</p><p>“Those three will be looked upon to fill most of what we lost last year,” Anderson said.</p><p><span class="-newsgate-paragraph-cci-subhead-"><b>Depth at forward</b></span></p><p>Sophomore forward Blake Maddalena also returns from the previous team, which finished 13-14 overall and was fourth in the seven-team OAA Red Division.</p><p>North Farmington-Harrison has added eight new forwards – senior Chad Hoffman, juniors Justin Bass, Max Flam and Alex Winkleman, sophomores Talon Brehmer, Jordan Pitts and Evan MacDonald and freshman Ben McColl.</p><p>Anderson has Cikra and Nathan skating on a line with MacDonald, a former travel player. Glasser will center the second line and have Maddalena and Bass on the wings.</p><p>“The third line is still a work in progress,” Anderson said, adding it will be for the first month of the season.</p><p>Three of the forwards are on the no-play list due to injury, ineligibility and the transfer rule.</p><p><span class="-newsgate-paragraph-cci-subhead-"><b>New leaders on ‘D’</b></span></p><p>The Flyers also lost their top two defensemen to graduation – all-stater Frank Zak, who contributed four goals and 30 assists, and Lucio D’Ascenzo.</p><p>Senior Tyler Magdich and junior Carlos Tobar played alongisde Zak and D’Ascenzo and gained a lot of experience. They’re the leaders now.</p><p>“It’s going to be a challenge, because Frank and Lucio were real good,” Anderson said. “It’ll be hard to replace them, too.</p><p>“Tyler had a great showing in our scrimmage. He stepped up and was an aggressive player; he scored a goal. If he can continue on that path, it will be really helpful to our team.”</p><p>Junior defenseman Ryan Kruger also returns and will have an expanded role this season.</p><p>Junior Thomas Zak was moved from forward to defense for more depth. Junior Jamie Dodd and freshman Daniel Bartlett are new defensemen.</p><p><span class="-newsgate-paragraph-cci-subhead-"><b>Go with veterans</b></span></p><p>“Magdich and Tobar are the most experienced and ready for any situation we put them in,” Anderson said. “We’ll split those two up and leave one of them out there on the ice.</p><p>“When it comes down to crucial times at the end of a period or end of a game, they’ll be out there on the same shift in those situations.”</p><p>Dodd had a good summer of growth and is ready to be one of the top four defensemen along with Kruger, according to Anderson.</p><p>“With the unknown in regard to our goal scoring, it’s going to be hard to match what other teams put up,” Anderson said.</p><p>“We’ll look to the defense to make the first pass out of the zone and clear it out. Hopefully, they’re up to it.”</p><p><span class="-newsgate-paragraph-cci-subhead-"><b>Goaltending duties</b></span></p><p>Brendan Dilloway did most of the goaltending, but senior Thomas Bacon and sophomore Colin Woods played 200-plus minutes apiece.</p><p>“I like our situation,” Anderson said. “They’re good hard-working kids. They’re going to push each other in practice and root for each other at the same time.</p><p>“Thomas is a firey kid; he anticipates getting out there, playing every game and doing well. As a senior, he’ll be a good example for Colin.”</p><p>The Flyers open the season in the Metro High School Invitational at Novi Ice Arena. They play Livonia Churchill at 5:30 p.m. Friday and South Lyon at 1:30 p.m. Saturday.</p><p>“With the turnover of players and talent we lost, it’s going to be hard to make up for that early on with all these new players as they adjust to the high school game,” Anderson said.</p><p>“With the league schedule and who we play, it’s not going to be easy. If we compete for the league championship and finish around the .500 mark overall, I think we’d take that at this point.”</p><div id="module-position-OlUUBRQKVZA" class="story-asset inline-share-tools-asset"><div class="inline-share-tools asset-inline-share-tools asset-inline-share-tools-bottom"><span class="inline-share-btn inline-share-btn-facebook" data-share-method="facebook" data-uotrack="InlineShareFacebookBtn"><span class="inline-share-count inline-share-count-facebook"></span><span class="inline-share-btn-label inline-share-btn-label-facebook">CONNECT</span></span><a class="inline-share-btn inline-share-btn-twitter" data-share-method="twitter" data-uotrack="InlineShareTwitterLink" data-popup-width="550" data-popup-height="450" href="https://twitter.com/intent/tweet?url=http%3A//www.hometownlife.com/story/sports/high-school/farmington/2015/11/17/flyers-need-find-scoring-punch/75912786/&amp;text=Flyers%20need%20to%20find%20more%20scoring%20punch&amp;via=hometownlife" target="_blank" data-popup=""><span class="inline-share-count inline-share-count-twitter"> 4 </span><span class="inline-share-btn-label inline-share-btn-label-twitter">TWEET</span></a><a class="inline-share-btn inline-share-btn-linkedin" data-share-method="linkedin" data-popup-width="600" data-popup-height="455" data-uotrack="UtilityBarFlyoutLinkedInLink" href="http://www.linkedin.com/shareArticle?url=http%3A//www.hometownlife.com/story/sports/high-school/farmington/2015/11/17/flyers-need-find-scoring-punch/75912786/&amp;mini=true" target="_blank" data-popup=""><span class="inline-share-count inline-share-count-linkedin"></span><span class="inline-share-btn-label inline-share-btn-label-linkedin">LINKEDIN</span></a><span class="inline-share-btn inline-share-btn-comments" data-share-method="comments" data-uotrack="InlineShareCommentsBtn"><span class="inline-share-count inline-share-count-comments"></span><span class="inline-share-btn-label inline-share-btn-label-comments">COMMENT</span></span><span class="inline-share-btn inline-share-btn-email" data-share-method="email" data-uotrack="InlineShareEmailBtn"><span class="inline-share-count inline-share-count-email"></span><span class="inline-share-btn-label inline-share-btn-label-email">EMAIL</span></span><span class="inline-share-btn inline-share-btn-more" data-share-method="facebook" data-uotrack="InlineShareMoreBtn"><span class="inline-share-count inline-share-count-more"></span><span class="inline-share-btn-label inline-share-btn-label-more">MORE</span></span></div></div><div class="article-print-url">Read or Share this story: http://www.hometownlife.com/story/sports/high-school/farmington/2015/11/17/flyers-need-find-scoring-punch/75912786/</div></div>`
var testExpectedText string = "On the eve of a new hockey season, the major issue facing Farmington Hills Unified coach Ken Anderson and his team is how to replace all the scoring that was lost. Austin Bottrell was an all-state player after the greatest offensive output by a Flyers player. He pumped in 34 goals and had 28 assists for 62 points. Danny Arnold and his 35 points was another graduation loss, and Joey Lajcaj (23 points) opted to play junior hockey in Colorado instead of returning for his senior season. “It’s going to be a struggle to make up that many points,” Anderson said. “It will have to be a team effort to get some of them. We’ll have to play a defensive style game and take care of the defensive zone first.” First of all, the Flyers will look to senior Jack Cikra, the leading scorer among the returning players. He had 17 goals and 24 points to rank third on the team. Seniors Brandon Glasser (15) and Andrew Nathan (18) will play prominent roles in helping to recharge the offense, too. Glasser had five goals and Nathan six. “Those three will be looked upon to fill most of what we lost last year,” Anderson said. Sophomore forward Blake Maddalena also returns from the previous team, which finished 13-14 overall and was fourth in the seven-team OAA Red Division. North Farmington-Harrison has added eight new forwards – senior Chad Hoffman, juniors Justin Bass, Max Flam and Alex Winkleman, sophomores Talon Brehmer, Jordan Pitts and Evan MacDonald and freshman Ben McColl. Anderson has Cikra and Nathan skating on a line with MacDonald, a former travel player. Glasser will center the second line and have Maddalena and Bass on the wings. “The third line is still a work in progress,” Anderson said, adding it will be for the first month of the season. Three of the forwards are on the no-play list due to injury, ineligibility and the transfer rule. The Flyers also lost their top two defensemen to graduation – all-stater Frank Zak, who contributed four goals and 30 assists, and Lucio D’Ascenzo. Senior Tyler Magdich and junior Carlos Tobar played alongisde Zak and D’Ascenzo and gained a lot of experience. They’re the leaders now. “It’s going to be a challenge, because Frank and Lucio were real good,” Anderson said. “It’ll be hard to replace them, too. “Tyler had a great showing in our scrimmage. He stepped up and was an aggressive player; he scored a goal. If he can continue on that path, it will be really helpful to our team.” Junior defenseman Ryan Kruger also returns and will have an expanded role this season. Junior Thomas Zak was moved from forward to defense for more depth. Junior Jamie Dodd and freshman Daniel Bartlett are new defensemen. “Magdich and Tobar are the most experienced and ready for any situation we put them in,” Anderson said. “We’ll split those two up and leave one of them out there on the ice. “When it comes down to crucial times at the end of a period or end of a game, they’ll be out there on the same shift in those situations.” Dodd had a good summer of growth and is ready to be one of the top four defensemen along with Kruger, according to Anderson. “With the unknown in regard to our goal scoring, it’s going to be hard to match what other teams put up,” Anderson said. “We’ll look to the defense to make the first pass out of the zone and clear it out. Hopefully, they’re up to it.” Brendan Dilloway did most of the goaltending, but senior Thomas Bacon and sophomore Colin Woods played 200-plus minutes apiece. “I like our situation,” Anderson said. “They’re good hard-working kids. They’re going to push each other in practice and root for each other at the same time. “Thomas is a firey kid; he anticipates getting out there, playing every game and doing well. As a senior, he’ll be a good example for Colin.” The Flyers open the season in the Metro High School Invitational at Novi Ice Arena. They play Livonia Churchill at 5:30 p.m. Friday and South Lyon at 1:30 p.m. Saturday. “With the turnover of players and talent we lost, it’s going to be hard to make up for that early on with all these new players as they adjust to the high school game,” Anderson said. “With the league schedule and who we play, it’s not going to be easy. If we compete for the league championship and finish around the .500 mark overall, I think we’d take that at this point.”"

func TestSubheadRemoval(t *testing.T) {
	html := bytes.NewBufferString(testHTML)
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		t.Fatal(err)
	}

	extractedBody := ExtractBodyFromDocument(doc, false, false)

	subhead := "Depth at forward"
	if strings.Contains(extractedBody.Text, subhead) {
		t.Fatalf("'%s' is a subhead, which should not appear in text body", subhead)
	}

	subhead = "New leaders on ‘D’"
	if strings.Contains(extractedBody.Text, subhead) {
		t.Fatalf("'%s' is a subhead, which should not appear in text body", subhead)
	}

	subhead = "Go with veterans"
	if strings.Contains(extractedBody.Text, subhead) {
		t.Fatalf("'%s' is a subhead, which should not appear in text body", subhead)
	}

	subhead = "Goaltending duties"
	if strings.Contains(extractedBody.Text, subhead) {
		t.Fatalf("'%s' is a subhead, which should not appear in text body", subhead)
	}

	actualText := strings.Join(strings.Fields(extractedBody.Text), " ")
	if actualText != testExpectedText {
		t.Fatal("Actiual text does not match expected text")
	}

}
func TestPhotoExtraction(t *testing.T) {
	t.Skip()
	doc, _ := goquery.NewDocument("http://www.freep.com/story/sports/nfl/lions/2015/11/20/jim-caldwell-hot-seat-lions/76103780/")
	photo := ExtractPhotoInfo(doc)
	fmt.Printf("%v", photo)
}

func TestTitleExtraction(t *testing.T) {
	urls := []string{
		"http://www.freep.com/picture-gallery/sports/high-school/2015/12/13/2015-free-press-dream-team-banquet/77265886/",
		"http://www.freep.com/story/sports/college/michigan-state/joe-rexrode/2015/12/15/michigan-state-football-basketball/77355566/",
		"http://www.detroitnews.com/picture-gallery/sports/nba/pistons/2015/11/30/pistons-116-rockets-105/76590382/",
	}

	for _, url := range urls {
		doc, err := goquery.NewDocument(url)
		if err != nil {
			t.Fatal(err)
		}

		title := ExtractTitleFromDocument(doc)
		if title == "" {
			t.Fatal("Failed to parse title from %s", url)
		}
	}
}